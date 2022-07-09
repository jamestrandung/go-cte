package die

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

type switchComputer interface {
	Switch(ctx context.Context, p any) (MasterPlan, error)
}

type bridgeComputer struct {
	se sideEffectComputer
	sw switchComputer
}

func (bc bridgeComputer) Compute(ctx context.Context, p any) (any, error) {
	if bc.se != nil {
		return struct{}{}, bc.se.Compute(ctx, p)
	}

	return bc.sw.Switch(ctx, p)
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
	outcome, _ := t.Outcome()
	return outcome.(V)
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
	if result, ok := o.(V); ok {
		return result
	}

	var tmp V
	return tmp
}
