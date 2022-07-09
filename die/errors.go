package die

import "errors"

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
	ErrPlanNotAnalyzed              = errors.New("plan must be analyzed before getting executed")
	ErrPlanMustUsePointerReceiver   = errors.New("the passed in plan must be a pointer")
)
