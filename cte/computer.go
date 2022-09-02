package cte

import (
    "context"
    "reflect"

    "github.com/jamestrandung/go-concurrency/async"
)

type ImpureComputer interface {
    Compute(ctx context.Context, p MasterPlan) (any, error)
}

type ImpureComputerWithLoadingData interface {
    LoadingComputer
    Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error)
}

type SideEffectComputer interface {
    Compute(ctx context.Context, p MasterPlan) error
}

type SideEffectComputerWithLoadingData interface {
    LoadingComputer
    Compute(ctx context.Context, p MasterPlan, data LoadingData) error
}

type SwitchComputer interface {
    Switch(ctx context.Context, p MasterPlan) (MasterPlan, error)
}

type SwitchComputerWithLoadingData interface {
    LoadingComputer
    Switch(ctx context.Context, p MasterPlan, data LoadingData) (MasterPlan, error)
}

type LoadingComputer interface {
    Load(ctx context.Context, p MasterPlan) (any, error)
}

type LoadingData struct {
    Data any
    Err  error
}

type toExecutePlan struct {
    mp MasterPlan
}

type loadingFn func(ctx context.Context, p MasterPlan) (any, error)

type delegatingComputer struct {
    loadingFn loadingFn
    computeFn func(ctx context.Context, p MasterPlan, data LoadingData) (any, error)
}

func (dc delegatingComputer) Load(ctx context.Context, p MasterPlan) (any, error) {
    if dc.loadingFn == nil {
        return nil, nil
    }

    return dc.loadingFn(ctx, p)
}

func (dc delegatingComputer) Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
    return dc.computeFn(ctx, p, data)
}

func newDelegatingComputer(rawComputer any) delegatingComputer {
    switch c := rawComputer.(type) {
    case ImpureComputerWithLoadingData:
        return delegatingComputer{
            loadingFn: func(ctx context.Context, p MasterPlan) (any, error) {
                return c.Load(ctx, p)
            },
            computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
                return c.Compute(ctx, p, data)
            },
        }
    case ImpureComputer:
        return delegatingComputer{
            computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
                return c.Compute(ctx, p)
            },
        }
    case SideEffectComputerWithLoadingData:
        return delegatingComputer{
            loadingFn: func(ctx context.Context, p MasterPlan) (any, error) {
                return c.Load(ctx, p)
            },
            computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
                return struct{}{}, c.Compute(ctx, p, data)
            },
        }
    case SideEffectComputer:
        return delegatingComputer{
            computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
                return struct{}{}, c.Compute(ctx, p)
            },
        }
    case SwitchComputerWithLoadingData:
        return delegatingComputer{
            loadingFn: func(ctx context.Context, p MasterPlan) (any, error) {
                return c.Load(ctx, p)
            },
            computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
                mp, err := c.Switch(ctx, p, data)

                return toExecutePlan{
                    mp: mp,
                }, err
            },
        }
    case SwitchComputer:
        return delegatingComputer{
            computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
                mp, err := c.Switch(ctx, p)

                return toExecutePlan{
                    mp: mp,
                }, err
            },
        }
    default:
        panic(ErrInvalidComputerType.Err(reflect.TypeOf(c)))
    }
}

type SideEffect struct{}

type SyncSideEffect struct{}

type Result struct {
    Task async.Task[any]
}

func newResult(t async.Task[any]) Result {
    return Result{
        Task: t,
    }
}

func Outcome[V any](t async.Task[any]) V {
    if t == nil {
        var tmp V
        return tmp
    }

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
