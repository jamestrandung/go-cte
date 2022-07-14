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
	"runtime/debug"
)

func (e Engine) IsAnalyzed(p plan) bool {
	_, ok := e.plans[extractFullNameFromValue(p)]
	return ok
}

func (e Engine) IsRegistered(v any) bool {
	fullName := extractFullNameFromValue(v)

	_, ok := e.computers[fullName]
	return ok
}

func (e Engine) IsExecutable(p masterPlan) (err error) {
	var verifyFn func(planName string)
	verifyFn = func(planName string) {
		ap := e.findAnalyzedPlan(planName)

		for _, component := range ap.components {
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("plan is not executable, %v", r)
						fmt.Println(string(debug.Stack()))
					}
				}()

				// If plan is not executable, 1 of the impureComputer will panic
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

func isValid(p plan) {
	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Pointer {
		panic(ErrPlanMustUsePointerReceiver)
	}

}
