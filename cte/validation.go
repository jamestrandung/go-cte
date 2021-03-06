package cte

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
)

func (e Engine) IsAnalyzed(p Plan) bool {
	_, ok := e.plans[extractFullNameFromValue(p)]
	return ok
}

func (e Engine) IsRegistered(v any) bool {
	fullName := extractFullNameFromValue(v)

	_, ok := e.computers[fullName]
	return ok
}

func (e Engine) IsExecutable(p MasterPlan) (err error) {
	var verifyFn func(planName string)
	verifyFn = func(planName string) {
		ap := e.findAnalyzedPlan(planName, reflect.ValueOf(p))

		for _, component := range ap.components {
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("plan is not executable, %v", r)
						fmt.Println(string(debug.Stack()))
					}
				}()

				// If plan is not executable, one of the computers will panic
				if c, ok := e.computers[component.id]; ok {
					c.Compute(context.Background(), p)
				}

				if _, ok := e.plans[component.id]; ok {
					verifyFn(component.id)
				}
			}()
		}
	}

	verifyFn(extractFullNameFromValue(p))

	return
}

func isValid(p Plan) {
	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Pointer {
		panic(ErrPlanMustUsePointerReceiver)
	}

}
