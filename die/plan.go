package die

import (
	"context"
)

type plan interface {
	IsSequential() bool
}

type MasterPlan interface {
	plan
	Execute(ctx context.Context) error
}

type pre interface {
	PreExecute(p any) error
}

type post interface {
	PostExecute(p any) error
}
