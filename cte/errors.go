package cte

import (
	"errors"
	"fmt"
)

type formatErr struct {
	format string
}

func makeFormatErr(format string) formatErr {
	return formatErr{
		format: format,
	}
}

func (e formatErr) Err(v ...any) error {
	return fmt.Errorf(e.format, v...)
}

var (
	// ErrPlanExecutionEndingEarly can be thrown actively by clients to end plan execution early.
	// For example, a value was retrieved from cache and thus, there's no point executing the algo
	// to calculate this value anymore. The engine will swallow this error, end execution and then
	// return a nil error to clients.
	//
	// Note: If the ending plan is nested inside another plan, the outer plan will still continue
	// as usual.
	ErrPlanExecutionEndingEarly = errors.New("plan execution ending early")
	// ErrRootPlanExecutionEndingEarly can be thrown actively by clients to end plan execution
	// early. For example, a value was retrieved from cache and thus, there's no point executing
	// the algo to calculate this value anymore. The engine will swallow this error, end execution
	// and then return a nil error to clients.
	//
	// Note: If the ending plan is nested inside another plan, the outer plan will also end.
	ErrRootPlanExecutionEndingEarly = errors.New("plan execution ending early from root")

	ErrPlanMustUsePointerReceiver = makeFormatErr("%v is using value receiver, all plans must be implemented using pointer receiver")
	ErrPlanNotAnalyzed            = makeFormatErr("%v has not been analyzed yet, call AnalyzePlan on it first")

	ErrInvalidComputerType = makeFormatErr("%v is not a computer")
	ErrNilMetadata         = makeFormatErr("metadata is missing for %v")
	ErrKeyMetaMissing      = makeFormatErr("key meta is missing for %v")
)
