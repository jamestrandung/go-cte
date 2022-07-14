package cte

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jamestrandung/go-concurrency/async"
	"golang.org/x/sync/errgroup"
)

var (
	planType       = reflect.TypeOf((*Plan)(nil)).Elem()
	preHookType    = reflect.TypeOf((*Pre)(nil)).Elem()
	postHookType   = reflect.TypeOf((*Post)(nil)).Elem()
	resultType     = reflect.TypeOf(Result{})
	syncResultType = reflect.TypeOf(SyncResult{})
)

type parsedComponent struct {
	id           string
	fieldIdx     int
	fieldType    reflect.Type
	isSyncResult bool
	requireSet   bool
}

type analyzedPlan struct {
	isSequential bool
	components   []parsedComponent
	preHooks     []Pre
	postHooks    []Post
}

type Engine struct {
	computers map[string]ImpureComputer
	plans     map[string]analyzedPlan
}

func NewEngine() Engine {
	return Engine{
		computers: make(map[string]ImpureComputer),
		plans:     make(map[string]analyzedPlan),
	}
}

func (e Engine) RegisterImpureComputer(v any, c ImpureComputer) {
	e.computers[extractFullNameFromValue(v)] = c
}

func (e Engine) RegisterSideEffectComputer(v any, c SideEffectComputer) {
	e.computers[extractFullNameFromValue(v)] = bridgeComputer{
		se: c,
	}
}

func (e Engine) RegisterSwitchComputer(v any, c SwitchComputer) {
	e.computers[extractFullNameFromValue(v)] = bridgeComputer{
		sw: c,
	}
}

func (e Engine) ConnectPreHook(p Plan, hooks ...Pre) {
	planName := extractFullNameFromValue(p)

	toUpdate := e.findExistingPlanOrCreate(planName)
	toUpdate.preHooks = append(toUpdate.preHooks, hooks...)

	e.plans[planName] = toUpdate
}

func (e Engine) ConnectPostHook(p Plan, hooks ...Post) {
	planName := extractFullNameFromValue(p)

	toUpdate := e.findExistingPlanOrCreate(planName)
	toUpdate.postHooks = append(toUpdate.postHooks, hooks...)

	e.plans[planName] = toUpdate
}

func (e Engine) AnalyzePlan(p Plan) string {
	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Pointer {
		panic(ErrPlanMustUsePointerReceiver)
	}

	val = val.Elem()

	var preHooks []Pre
	var postHooks []Post

	components := make([]parsedComponent, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i).Type
		fieldPointerType := reflect.PointerTo(fieldType)

		// Hook types might be embedded in a parent plan struct. Hence, we need to check if the type
		// is a hook but not a plan so that we don't register a plan as a hook.
		typeAndPointerTypeIsNotPlanType := !fieldType.Implements(planType) && !fieldPointerType.Implements(planType)

		// Hooks might be implemented with value or pointer receivers.
		isPreHookType := fieldType.Implements(preHookType) || fieldPointerType.Implements(preHookType)
		isPostHookType := fieldType.Implements(postHookType) || fieldPointerType.Implements(postHookType)

		if typeAndPointerTypeIsNotPlanType && isPreHookType {
			preHook := reflect.New(fieldType).Interface().(Pre)
			preHooks = append(preHooks, preHook)

			continue
		}

		if typeAndPointerTypeIsNotPlanType && isPostHookType {
			postHook := reflect.New(fieldType).Interface().(Post)
			postHooks = append(postHooks, postHook)

			continue
		}

		componentID := extractFullNameFromType(fieldType)

		component := func() parsedComponent {
			if fieldType.ConvertibleTo(resultType) {
				// Both sequential & parallel plans can contain Result fields
				return parsedComponent{
					id:         componentID,
					fieldIdx:   i,
					fieldType:  fieldType,
					requireSet: true,
				}
			}

			if fieldType.ConvertibleTo(syncResultType) {
				if !p.IsSequential() {
					panic(fmt.Errorf("parallel plan cannot contain SyncResult field: %s", extractShortName(componentID)))
				}

				return parsedComponent{
					id:           componentID,
					fieldIdx:     i,
					fieldType:    fieldType,
					isSyncResult: true,
					requireSet:   true,
				}
			}

			return parsedComponent{
				id: componentID,
			}
		}()

		components[i] = component
	}

	planName := extractFullNameFromValue(p)

	toUpdate := e.findExistingPlanOrCreate(planName)
	toUpdate.isSequential = p.IsSequential()
	toUpdate.components = components
	toUpdate.preHooks = append(toUpdate.preHooks, preHooks...)
	toUpdate.postHooks = append(toUpdate.postHooks, postHooks...)

	e.plans[planName] = toUpdate

	return planName
}

func (e Engine) ExecuteMasterPlan(ctx context.Context, planName string, p MasterPlan) error {
	// Plan implementations always use pointer receivers.
	// Should be safe to extract value.
	planValue := reflect.ValueOf(p).Elem()

	if err := e.doExecute(ctx, planName, p, planValue, p.IsSequential()); err != nil {
		return swallowErrPlanExecutionEndingEarly(err)
	}

	return nil
}

func (e Engine) doExecute(ctx context.Context, planName string, p MasterPlan, curPlanValue reflect.Value, isSequential bool) error {
	ap := e.findAnalyzedPlan(planName)

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

func (e Engine) doExecuteSync(ctx context.Context, p MasterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
	for idx, component := range components {
		if c, ok := e.computers[component.id]; ok {
			task := async.NewTask(
				func(taskCtx context.Context) (any, error) {
					result, err := c.Compute(taskCtx, p)
					if tep, ok := result.(toExecutePlan); err == nil && ok {
						return tep.mp, tep.mp.Execute(taskCtx)
					}

					return result, err
				},
			)

			result, err := task.RunSync(ctx).Outcome()
			if err != nil {
				return err
			}

			// Register Result/SyncResult in a sequential plan's field
			if component.requireSet {
				field := curPlanValue.Field(component.fieldIdx)
				casted := func() reflect.Value {
					if component.isSyncResult {
						return reflect.ValueOf(newSyncResult(result)).Convert(component.fieldType)
					}

					return reflect.ValueOf(newResult(task)).Convert(component.fieldType)
				}()

				field.Set(casted)
			}

			continue
		}

		// Nested plan gets executed synchronously
		if ap, ok := e.plans[component.id]; ok {
			// Nested plan is always a value, never a pointer. Hence, no need to call Elem().
			nestedPlanValue := curPlanValue.Field(idx)

			err := e.doExecute(ctx, component.id, p, nestedPlanValue, ap.isSequential)
			if err != nil && err != ErrPlanExecutionEndingEarly {
				return err
			}
		}
	}

	return nil
}

func (e Engine) doExecuteAsync(ctx context.Context, p MasterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
	tasks := make([]async.SilentTask, 0, len(components))
	for idx, component := range components {
		componentID := component.id

		if c, ok := e.computers[componentID]; ok {
			task := async.NewTask(
				func(taskCtx context.Context) (any, error) {
					result, err := c.Compute(taskCtx, p)
					if tep, ok := result.(toExecutePlan); err == nil && ok {
						return tep.mp, tep.mp.Execute(taskCtx)
					}

					return result, err
				},
			)

			tasks = append(tasks, task)

			// Register Result in a parallel plan's field
			if component.requireSet {
				field := curPlanValue.Field(component.fieldIdx)
				casted := reflect.ValueOf(newResult(task)).Convert(component.fieldType)

				field.Set(casted)
			}

			continue
		}

		// Nested plan gets executed asynchronously by wrapping it inside a task
		if ap, ok := e.plans[componentID]; ok {
			// Nested plan is always a value, never a pointer. Hence, no need to call Elem().
			nestedPlanValue := curPlanValue.Field(idx)

			task := async.NewSilentTask(
				func(taskCtx context.Context) error {
					err := e.doExecute(taskCtx, component.id, p, nestedPlanValue, ap.isSequential)
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

func (e Engine) findAnalyzedPlan(planName string) analyzedPlan {
	ap, ok := e.plans[planName]
	if !ok || len(ap.components) == 0 {
		panic(ErrPlanNotAnalyzed)
	}

	return ap
}
