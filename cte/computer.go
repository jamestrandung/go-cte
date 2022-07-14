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

	"github.com/jamestrandung/go-concurrency/async"
)

type impureComputer interface {
	Compute(ctx context.Context, p any) (any, error)
}

type sideEffectComputer interface {
	Compute(ctx context.Context, p any) error
}

type bridgeComputer struct {
	sc sideEffectComputer
}

func (bc bridgeComputer) Compute(ctx context.Context, p any) (any, error) {
	return struct{}{}, bc.sc.Compute(ctx, p)
}

type SideEffectKey struct{}

type Result struct {
	Task async.Task[any]
}

func newAsyncResult(t async.Task[any]) Result {
	return Result{
		Task: t,
	}
}

func Outcome[V any](t async.Task[any]) V {
	result, _ := t.Outcome()
	return result.(V)
}

type SyncResult struct {
	Outcome any
}

func newSyncResult(o any) SyncResult {
	return SyncResult{
		Outcome: o,
	}
}

func Cast[V any](o any) V {
	return o.(V)
}
