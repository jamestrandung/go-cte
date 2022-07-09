package die

import (
	"context"

	"github.com/jamestrandung/go-concurrency/async"
)

type ImpureComputer interface {
	Compute(ctx context.Context, p any) (any, error)
}

type SideEffectComputer interface {
	Compute(ctx context.Context, p any) error
}

type SwitchComputer interface {
	Switch(ctx context.Context, p any) (MasterPlan, error)
}

type toExecutePlan struct {
	mp MasterPlan
}

type bridgeComputer struct {
	se SideEffectComputer
	sw SwitchComputer
}

func (bc bridgeComputer) Compute(ctx context.Context, p any) (any, error) {
	if bc.se != nil {
		return struct{}{}, bc.se.Compute(ctx, p)
	}

	mp, err := bc.sw.Switch(ctx, p)

	return toExecutePlan{
		mp: mp,
	}, err
}

type SideEffectKey struct{}

type Result struct {
	Task async.Task[any]
}

func newResult(t async.Task[any]) Result {
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
