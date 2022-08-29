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

var emptyLoadingData = LoadingData{}

type LoadingData struct {
	Data any
	Err  error
}

type toExecutePlan struct {
	mp MasterPlan
}

type bridgeComputer struct {
	computeFn func(ctx context.Context, p MasterPlan, data LoadingData) (any, error)
}

func (bc bridgeComputer) Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
	return bc.computeFn(ctx, p, data)
}

type computerWrapper struct {
	LoadingComputer
	bridgeComputer
}

func newComputerWrapper(rawComputer any) computerWrapper {
	switch c := rawComputer.(type) {
	case ImpureComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer: c,
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
					return c.Compute(ctx, p, data)
				},
			},
		}
	case ImpureComputer:
		return computerWrapper{
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
					return c.Compute(ctx, p)
				},
			},
		}
	case SideEffectComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer: c,
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
					return struct{}{}, c.Compute(ctx, p, data)
				},
			},
		}
	case SideEffectComputer:
		return computerWrapper{
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
					return struct{}{}, c.Compute(ctx, p)
				},
			},
		}
	case SwitchComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer: c,
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
					mp, err := c.Switch(ctx, p, data)

					return toExecutePlan{
						mp: mp,
					}, err
				},
			},
		}
	case SwitchComputer:
		return computerWrapper{
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
					mp, err := c.Switch(ctx, p)

					return toExecutePlan{
						mp: mp,
					}, err
				},
			},
		}
	default:
		panic(ErrInvalidComputerType.Err(reflect.TypeOf(c)))
	}
}

func (w computerWrapper) Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
	return w.bridgeComputer.Compute(ctx, p, data)
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
