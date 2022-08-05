package cte

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jamestrandung/go-concurrency/async"
	"golang.org/x/sync/errgroup"
)

var (
	planType             = reflect.TypeOf((*Plan)(nil)).Elem()
	preHookType          = reflect.TypeOf((*Pre)(nil)).Elem()
	postHookType         = reflect.TypeOf((*Post)(nil)).Elem()
	metadataProviderType = reflect.TypeOf((*MetadataProvider)(nil)).Elem()
	resultType           = reflect.TypeOf(Result{})
	syncResultType       = reflect.TypeOf(SyncResult{})
)

type registeredComputer struct {
	computer ImpureComputer
	metadata map[metaType]reflect.Type
}

type parsedComponent struct {
	id            string
	fieldIdx      int
	fieldType     reflect.Type
	isSyncResult  bool
	requireSet    bool
	isPointerType bool
}

type analyzedPlan struct {
	isSequential bool
	components   []parsedComponent
	preHooks     []Pre
	postHooks    []Post
}

type Engine struct {
	computers map[string]registeredComputer
	plans     map[string]analyzedPlan
}

func NewEngine() Engine {
	return Engine{
		computers: make(map[string]registeredComputer),
		plans:     make(map[string]analyzedPlan),
	}
}

func (e Engine) AnalyzePlan(p Plan) {
	planName := extractFullNameFromValue(p)
	if _, ok := e.plans[planName]; ok {
		return
	}

	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Pointer {
		panic(ErrPlanMustUsePointerReceiver.Err(reflect.TypeOf(p)))
	}

	val = val.Elem()

	var preHooks []Pre
	var postHooks []Post
	var components []parsedComponent

	for i := 0; i < val.NumField(); i++ {
		isPointerType, fieldType, fieldPointerType := extractTypes(val.Type().Field(i))

		processedAsHook := func() bool {
			// Hook types might be embedded in a parent plan struct. Hence, we need to check if the type
			// is a hook but not a plan so that we don't register a plan as a hook.
			typeAndPointerTypeIsNotPlanType := !fieldType.Implements(planType) && !fieldPointerType.Implements(planType)

			// Hooks might be implemented with value or pointer receivers.
			isPreHookType := fieldType.Implements(preHookType) || fieldPointerType.Implements(preHookType)
			isPostHookType := fieldType.Implements(postHookType) || fieldPointerType.Implements(postHookType)

			if typeAndPointerTypeIsNotPlanType && isPreHookType {
				// Call to Interface() returns a pointer value which is acceptable for
				// both scenarios where fieldType uses pointer or value receiver to
				// implement an interface
				preHook := reflect.New(fieldType).Interface().(Pre)
				preHooks = append(preHooks, preHook)

				return true
			}

			if typeAndPointerTypeIsNotPlanType && isPostHookType {
				// Call to Interface() returns a pointer value which is acceptable for
				// both scenarios where fieldType uses pointer or value receiver to
				// implement an interface
				postHook := reflect.New(fieldType).Interface().(Post)
				postHooks = append(postHooks, postHook)

				return true
			}

			return false
		}()

		if processedAsHook {
			continue
		}

		componentID := extractFullNameFromType(fieldType)

		// Dynamically register computers that are actually used in a plan
		func() {
			isMetadataProviderType := fieldType.Implements(metadataProviderType) || fieldPointerType.Implements(metadataProviderType)
			if !isMetadataProviderType {
				return
			}

			mp := reflect.New(fieldType).Interface().(MetadataProvider)
			e.registerComputer(mp)
		}()

		// Dynamically analyze nested plans
		func() {
			isPlanType := fieldPointerType.Implements(planType)

			if !isPlanType {
				return
			}

			if isPointerType {
				panic(ErrNestedPlanCannotBePointer.Err(val.Type(), fieldType))
			}

			e.AnalyzePlan(reflect.New(fieldType).Interface().(Plan))
		}()

		component := func() parsedComponent {
			if fieldType.ConvertibleTo(resultType) {
				// Both sequential & parallel plans can contain Result fields
				return parsedComponent{
					id:            componentID,
					fieldIdx:      i,
					fieldType:     fieldType,
					requireSet:    true,
					isPointerType: isPointerType,
				}
			}

			if fieldType.ConvertibleTo(syncResultType) {
				if !p.IsSequentialCTEPlan() {
					panic(fmt.Errorf("parallel plan cannot contain SyncResult field: %s", extractShortName(componentID)))
				}

				return parsedComponent{
					id:            componentID,
					fieldIdx:      i,
					fieldType:     fieldType,
					isSyncResult:  true,
					requireSet:    true,
					isPointerType: isPointerType,
				}
			}

			return parsedComponent{
				id:       componentID,
				fieldIdx: i,
			}
		}()

		components = append(components, component)
	}

	e.plans[planName] = analyzedPlan{
		isSequential: p.IsSequentialCTEPlan(),
		components:   components,
		preHooks:     preHooks,
		postHooks:    postHooks,
	}
}

func (e Engine) registerComputer(mp MetadataProvider) {
	computerID := extractFullNameFromValue(mp)
	if _, ok := e.computers[computerID]; ok {
		return
	}

	metadata := extractMetadata(mp)

	computerType := func() reflect.Type {
		cType, ok := metadata[metaTypeComputer]
		if !ok {
			panic(ErrComputerMetaMissing.Err(reflect.TypeOf(mp)))
		}

		if cType.Kind() == reflect.Pointer {
			return cType.Elem()
		}

		return cType
	}()

	computer := reflect.New(computerType).Interface()

	switch c := computer.(type) {
	case ImpureComputer:
		e.computers[computerID] = registeredComputer{
			computer: c,
			metadata: metadata,
		}
	case SideEffectComputer:
		e.computers[computerID] = registeredComputer{
			computer: bridgeComputer{
				se: c,
			},
			metadata: metadata,
		}
	case SwitchComputer:
		e.computers[computerID] = registeredComputer{
			computer: bridgeComputer{
				sw: c,
			},
			metadata: metadata,
		}
	default:
		panic(ErrInvalidComputerType.Err(computerType))
	}
}

func (e Engine) ExecuteMasterPlan(ctx context.Context, p MasterPlan) error {
	// Plan implementations always use pointer receivers.
	// Should be safe to extract value.
	planValue := reflect.ValueOf(p).Elem()

	planName := extractFullNameFromType(planValue.Type())
	if err := e.doExecutePlan(ctx, planName, p, planValue, p.IsSequentialCTEPlan()); err != nil {
		return swallowErrPlanExecutionEndingEarly(err)
	}

	return nil
}

func (e Engine) doExecutePlan(ctx context.Context, planName string, p MasterPlan, curPlanValue reflect.Value, isSequential bool) error {
	ap := e.findAnalyzedPlan(planName, curPlanValue)

	for _, h := range ap.preHooks {
		if err := h.PreExecute(p); err != nil {
			return err
		}
	}

	err := func() error {
		if isSequential {
			return e.doExecuteSync(ctx, p, curPlanValue, ap.components)
		}

		return e.doExecuteAsync(ctx, p, curPlanValue, ap.components)
	}()

	if err != nil {
		return err
	}

	for _, h := range ap.postHooks {
		if err := h.PostExecute(p); err != nil {
			return err
		}
	}

	return nil
}

func (e Engine) doExecuteComputer(ctx context.Context, c ImpureComputer, p MasterPlan) (any, error) {
	result, err := c.Compute(ctx, p)
	if tep, ok := result.(toExecutePlan); ok {
		if err != nil {
			return tep.mp, err
		}

		return tep.mp, tep.mp.Execute(ctx)
	}

	return result, err
}

func (e Engine) doExecuteSync(ctx context.Context, p MasterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
	for _, component := range components {
		if c, ok := e.computers[component.id]; ok {
			task := async.NewTask(
				func(taskCtx context.Context) (any, error) {
					return e.doExecuteComputer(taskCtx, c.computer, p)
				},
			)

			result, err := task.RunSync(ctx).Outcome()

			// Register Result/SyncResult in a sequential plan's field
			if component.requireSet {
				field := curPlanValue.Field(component.fieldIdx)

				casted := func() reflect.Value {
					if component.isSyncResult {
						return reflect.ValueOf(newSyncResult(result)).Convert(component.fieldType)
					}

					return reflect.ValueOf(newResult(task)).Convert(component.fieldType)
				}()

				resultTakingIntoAccountPointerType := func() reflect.Value {
					if !component.isPointerType {
						return casted
					}

					rv := reflect.New(component.fieldType)
					rv.Elem().Set(casted)

					return rv
				}()

				field.Set(resultTakingIntoAccountPointerType)
			}

			if err != nil {
				return err
			}

			continue
		}

		// Nested plan gets executed synchronously
		if ap, ok := e.plans[component.id]; ok {
			// Nested plan is always a value, never a pointer. Hence, no need to call Elem().
			nestedPlanValue := curPlanValue.Field(component.fieldIdx)

			err := e.doExecutePlan(ctx, component.id, p, nestedPlanValue, ap.isSequential)
			if err != nil && err != ErrPlanExecutionEndingEarly {
				return err
			}
		}
	}

	return nil
}

func (e Engine) doExecuteAsync(ctx context.Context, p MasterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
	tasks := make([]async.SilentTask, 0, len(components))
	for _, component := range components {
		componentID := component.id

		if c, ok := e.computers[componentID]; ok {
			task := async.NewTask(
				func(taskCtx context.Context) (any, error) {
					return e.doExecuteComputer(taskCtx, c.computer, p)
				},
			)

			tasks = append(tasks, task)

			// Register Result in a parallel plan's field
			if component.requireSet {
				field := curPlanValue.Field(component.fieldIdx)

				resultTakingIntoAccountPointerType := func() reflect.Value {
					casted := reflect.ValueOf(newResult(task)).Convert(component.fieldType)
					if !component.isPointerType {
						return casted
					}

					rv := reflect.New(component.fieldType)
					rv.Elem().Set(casted)

					return rv
				}()

				field.Set(resultTakingIntoAccountPointerType)
			}

			continue
		}

		// Nested plan gets executed asynchronously by wrapping it inside a task
		if ap, ok := e.plans[componentID]; ok {
			// Nested plan is always a value, never a pointer. Hence, no need to call Elem().
			nestedPlanValue := curPlanValue.Field(component.fieldIdx)

			task := async.NewSilentTask(
				func(taskCtx context.Context) error {
					err := e.doExecutePlan(taskCtx, componentID, p, nestedPlanValue, ap.isSequential)
					if err != nil && err != ErrPlanExecutionEndingEarly {
						return err
					}

					return nil
				},
			)

			tasks = append(tasks, task)
		}
	}

	g, groupCtx := errgroup.WithContext(ctx)
	for _, task := range tasks {
		t := task
		g.Go(
			func() error {
				return t.ExecuteSync(groupCtx).Error()
			},
		)
	}

	return g.Wait()
}

func (e Engine) findExistingPlanOrCreate(planName string) analyzedPlan {
	if existing, ok := e.plans[planName]; ok {
		return existing
	}

	return analyzedPlan{}
}

func (e Engine) findAnalyzedPlan(planName string, curPlanValue reflect.Value) analyzedPlan {
	if len(planName) == 0 {
		panic(ErrPlanNotAnalyzed.Err(curPlanValue.Type()))
	}

	ap, ok := e.plans[planName]
	if !ok || len(ap.components) == 0 {
		panic(ErrPlanNotAnalyzed.Err(planName))
	}

	return ap
}
