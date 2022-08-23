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

type computerWrapper struct {
	LoadingComputer

	ImpureComputer
	ImpureComputerWithLoadingData
	SideEffectComputer
	SideEffectComputerWithLoadingData
	SwitchComputer
	SwitchComputerWithLoadingData
}

func newComputerWrapper(rawComputer any) computerWrapper {
	switch c := rawComputer.(type) {
	case ImpureComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer:               c,
			ImpureComputerWithLoadingData: c,
		}
	case ImpureComputer:
		return computerWrapper{
			ImpureComputer: c,
		}
	case SideEffectComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer:                   c,
			SideEffectComputerWithLoadingData: c,
		}
	case SideEffectComputer:
		return computerWrapper{
			SideEffectComputer: c,
		}
	case SwitchComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer:               c,
			SwitchComputerWithLoadingData: c,
		}
	case SwitchComputer:
		return computerWrapper{
			SwitchComputer: c,
		}
	default:
		panic(ErrInvalidComputerType.Err(reflect.TypeOf(c)))
	}
}

func (w computerWrapper) Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
	if w.ImpureComputer != nil {
		return w.ImpureComputer.Compute(ctx, p)
	}

	if w.SideEffectComputerWithLoadingData != nil {
		return struct{}{}, w.SideEffectComputerWithLoadingData.Compute(ctx, p, data)
	}

	if w.SideEffectComputer != nil {
		return struct{}{}, w.SideEffectComputer.Compute(ctx, p)
	}

	mp, err := func() (MasterPlan, error) {
		if w.SwitchComputerWithLoadingData != nil {
			return w.SwitchComputerWithLoadingData.Switch(ctx, p, data)
		}

		return w.SwitchComputer.Switch(ctx, p)
	}()

	return toExecutePlan{
		mp: mp,
	}, err
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
