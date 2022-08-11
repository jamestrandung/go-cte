package cte

import (
	"context"
	"reflect"

	"github.com/jamestrandung/go-concurrency/async"
)

type ImpureComputer interface {
	LoadingComputer
	Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error)
}

type ImpureComputerWithoutLoadingData interface {
	Compute(ctx context.Context, p MasterPlan) (any, error)
}

type SideEffectComputer interface {
	LoadingComputer
	Compute(ctx context.Context, p MasterPlan, data LoadingData) error
}

type SideEffectComputerWithoutLoadingData interface {
	Compute(ctx context.Context, p MasterPlan) error
}

type SwitchComputer interface {
	LoadingComputer
	Switch(ctx context.Context, p MasterPlan, data LoadingData) (MasterPlan, error)
}

type SwitchComputerWithoutLoadingData interface {
	Switch(ctx context.Context, p MasterPlan) (MasterPlan, error)
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

type computerParts struct {
	LoadingComputer

	ImpureComputer
	ImpureComputerWithoutLoadingData
	SideEffectComputer
	SideEffectComputerWithoutLoadingData
	SwitchComputer
	SwitchComputerWithoutLoadingData
}

func newComputerParts(rawComputer any) computerParts {
	switch c := rawComputer.(type) {
	case ImpureComputer:
		return computerParts{
			LoadingComputer: c,
			ImpureComputer:  c,
		}
	case ImpureComputerWithoutLoadingData:
		return computerParts{
			ImpureComputerWithoutLoadingData: c,
		}
	case SideEffectComputer:
		return computerParts{
			LoadingComputer:    c,
			SideEffectComputer: c,
		}
	case SideEffectComputerWithoutLoadingData:
		return computerParts{
			SideEffectComputerWithoutLoadingData: c,
		}
	case SwitchComputer:
		return computerParts{
			LoadingComputer: c,
			SwitchComputer:  c,
		}
	case SwitchComputerWithoutLoadingData:
		return computerParts{
			SwitchComputerWithoutLoadingData: c,
		}
	default:
		panic(ErrInvalidComputerType.Err(reflect.TypeOf(c)))
	}
}

func (cp computerParts) Compute(ctx context.Context, p MasterPlan, data LoadingData) (any, error) {
	if cp.ImpureComputerWithoutLoadingData != nil {
		return cp.ImpureComputerWithoutLoadingData.Compute(ctx, p)
	}

	if cp.SideEffectComputer != nil {
		return struct{}{}, cp.SideEffectComputer.Compute(ctx, p, data)
	}

	if cp.SideEffectComputerWithoutLoadingData != nil {
		return struct{}{}, cp.SideEffectComputerWithoutLoadingData.Compute(ctx, p)
	}

	mp, err := func() (MasterPlan, error) {
		if cp.SwitchComputer != nil {
			return cp.SwitchComputer.Switch(ctx, p, data)
		}

		return cp.SwitchComputerWithoutLoadingData.Switch(ctx, p)
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
