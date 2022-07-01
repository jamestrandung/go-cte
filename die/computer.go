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

type ComputerKey struct{}

type AsyncResult struct {
	Task async.Task[any]
}

func newAsyncResult(t async.Task[any]) AsyncResult {
	return AsyncResult{
		Task: t,
	}
}

func ExtractAsync[V any](t async.Task[any]) V {
	result, _ := t.Outcome()
	return result.(V)
}

type SyncResult struct {
	outcome any
}

func newSyncResult(o any) SyncResult {
	return SyncResult{
		outcome: o,
	}
}

func ExtractSync[V any](r SyncResult) V {
	return r.outcome.(V)
}
