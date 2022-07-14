// Copyright (c) 2012-2022 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package cte

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jamestrandung/go-concurrency/async"
	"golang.org/x/sync/errgroup"
)

var (
	planType       = reflect.TypeOf((*plan)(nil)).Elem()
	preHookType    = reflect.TypeOf((*pre)(nil)).Elem()
	postHookType   = reflect.TypeOf((*post)(nil)).Elem()
	resultType     = reflect.TypeOf(Result{})
	syncResultType = reflect.TypeOf(SyncResult{})
)

type parsedComponent struct {
	id     string
	setter *reflect.Method
}

type analyzedPlan struct {
	isSequential bool
	components   []parsedComponent
	preHooks     []pre
	postHooks    []post
}

type Engine struct {
	computers map[string]impureComputer
	plans     map[string]analyzedPlan
}

func NewEngine() Engine {
	return Engine{
		computers: make(map[string]impureComputer),
		plans:     make(map[string]analyzedPlan),
	}
}

func (e Engine) RegisterImpureComputer(v any, c impureComputer) {
	e.computers[extractFullNameFromValue(v)] = c
}

func (e Engine) RegisterSideEffectComputer(v any, sc sideEffectComputer) {
	e.computers[extractFullNameFromValue(v)] = bridgeComputer{
		sc: sc,
	}
}

func (e Engine) ConnectPreHook(p plan, hooks ...pre) {
	planName := extractFullNameFromValue(p)

	toUpdate := e.findExistingPlanOrCreate(planName)
	toUpdate.preHooks = append(toUpdate.preHooks, hooks...)

	e.plans[planName] = toUpdate
}

func (e Engine) ConnectPostHook(p plan, hooks ...post) {
	planName := extractFullNameFromValue(p)

	toUpdate := e.findExistingPlanOrCreate(planName)
	toUpdate.postHooks = append(toUpdate.postHooks, hooks...)

	e.plans[planName] = toUpdate
}

func (e Engine) AnalyzePlan(p plan) string {
	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Pointer {
		panic(ErrPlanMustUsePointerReceiver)
	}

	val = val.Elem()
	pType := reflect.ValueOf(p).Type()

	var preHooks []pre
	var postHooks []post

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
			preHook := reflect.New(fieldType).Interface().(pre)
			preHooks = append(preHooks, preHook)

			continue
		}

		if typeAndPointerTypeIsNotPlanType && isPostHookType {
			postHook := reflect.New(fieldType).Interface().(post)
			postHooks = append(postHooks, postHook)

			continue
		}

		componentID := extractFullNameFromType(fieldType)

		component := func() parsedComponent {
			if fieldType.ConvertibleTo(resultType) {
				// Both sequential & parallel plans can contain Result fields
				if setter, ok := pType.MethodByName("Set" + extractShortName(componentID)); ok {
					return parsedComponent{
						id:     componentID,
						setter: &setter,
					}
				}

				panic(fmt.Errorf("plan must have setter for Result field: %s", extractShortName(componentID)))
			}

			if fieldType.ConvertibleTo(syncResultType) {
				if !p.IsSequential() {
					panic(fmt.Errorf("parallel plan cannot contain SyncResult field: %s", extractShortName(componentID)))
				}

				if setter, ok := pType.MethodByName("Set" + extractShortName(componentID)); ok {
					return parsedComponent{
						id:     componentID,
						setter: &setter,
					}
				}

				panic(fmt.Errorf("sequential plan must have setter for SyncResult field: %s", extractShortName(componentID)))
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

func (e Engine) ExecuteMasterPlan(ctx context.Context, planName string, p masterPlan) error {
	// Plan implementations always use pointer receivers.
	// Should be safe to extract value.
	planValue := reflect.ValueOf(p).Elem()

	if err := e.doExecute(ctx, planName, p, planValue, p.IsSequential()); err != nil {
		return swallowErrPlanExecutionEndingEarly(err)
	}

	return nil
}

func (e Engine) doExecute(ctx context.Context, planName string, p masterPlan, curPlanValue reflect.Value, isSequential bool) error {
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

func (e Engine) doExecuteSync(ctx context.Context, p masterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
	for idx, component := range components {
		if c, ok := e.computers[component.id]; ok {
			task := async.NewTask(
				func(taskCtx context.Context) (any, error) {
					return c.Compute(taskCtx, p)
				},
			)

			result, err := task.RunSync(ctx).Outcome()
			if err != nil {
				return err
			}

			// Register SyncResult in a sequential plan's field
			if component.setter != nil {
				// Plan implementations always use pointer receivers.
				// Need to extract pointer in order to call setters.
				pointer := curPlanValue.Addr()

				component.setter.Func.Call([]reflect.Value{pointer, reflect.ValueOf(newSyncResult(result))})
			}

			continue
		}

		// Nested plan gets executed synchronously
		if ap, ok := e.plans[component.id]; ok {
			// Nested plan is always a value, never a pointer. Hence, no need to call Elem().
			nestedPlanValue := curPlanValue.Field(idx)

			if err := e.doExecute(ctx, component.id, p, nestedPlanValue, ap.isSequential); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e Engine) doExecuteAsync(ctx context.Context, p masterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
	tasks := make([]async.SilentTask, 0, len(components))
	for idx, component := range components {
		componentID := component.id

		if c, ok := e.computers[componentID]; ok {
			task := async.NewTask(
				func(taskCtx context.Context) (any, error) {
					return c.Compute(taskCtx, p)
				},
			)

			tasks = append(tasks, task)

			// Register Result in a parallel plan's field
			if component.setter != nil {
				// Plan implementations always use pointer receivers.
				// Need to extract pointer in order to call setters.
				pointer := curPlanValue.Addr()

				component.setter.Func.Call([]reflect.Value{pointer, reflect.ValueOf(newAsyncResult(task))})
			}

			continue
		}

		// Nested plan gets executed asynchronously by wrapping it inside a task
		if ap, ok := e.plans[componentID]; ok {
			// Nested plan is always a value, never a pointer. Hence, no need to call Elem().
			nestedPlanValue := curPlanValue.Field(idx)

			task := async.NewSilentTask(
				func(taskCtx context.Context) error {
					return e.doExecute(taskCtx, componentID, p, nestedPlanValue, ap.isSequential)
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
