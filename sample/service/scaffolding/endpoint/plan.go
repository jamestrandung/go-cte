package endpoint

import (
	"context"
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/components/quote"
	"github.com/jamestrandung/go-cte/sample/service/components/streaming"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/loading"
)

type SequentialPlan struct {
	Input
	Dependencies
	loading.ParallelPlan
	quote.CalculatedCost
	streaming.CostStreaming
}

func NewPlan(r Input, d Dependencies) *SequentialPlan {
	return &SequentialPlan{
		Input:        r,
		Dependencies: d,
	}
}

func (p *SequentialPlan) IsSequential() bool {
	return true
}

func (p *SequentialPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, planName, p)
}
