package cte

import (
	"context"
)

type Plan interface {
	IsSequential() bool
}

type MasterPlan interface {
	Plan
	Execute(ctx context.Context) error
}

type Pre interface {
	PreExecute(p any) error
}

type Post interface {
	PostExecute(p any) error
}
