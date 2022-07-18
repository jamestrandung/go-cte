package cte

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
    // We only need to take care of situations where clients return "nil, err"
    // in their computers. Other than that, if clients return a non-nil value,
    // they must make sure the return type is as expected.
    outcome, _ := t.Outcome()
    if outcome != nil {
        return outcome.(V)
    }

    var tmp V
    return tmp
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
    // We only need to take care of situations where clients return "nil, err"
    // in their computers. Other than that, if clients return a non-nil value,
    // they must make sure the return type is as expected.
    if o != nil {
        return o.(V)
    }

    var tmp V
    return tmp
}
