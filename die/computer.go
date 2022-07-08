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
