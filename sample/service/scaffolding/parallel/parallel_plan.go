package parallel

import (
	"context"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte/sample/service/components/travelcost"
	"github.com/jamestrandung/go-cte/sample/service/components/travelplan"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/sequential"
)

type ParallelPlan struct {
	Request
	Dependencies
	costconfigs.CostConfigs
	travelplan.TravelPlan
	travelcost.TravelCost
	sequential.SequentialPlan
}

func NewPlan(r Request, d Dependencies) *ParallelPlan {
	return &ParallelPlan{
		Request:      r,
		Dependencies: d,
	}
}

func (p *ParallelPlan) IsSequential() bool {
	return false
}

func (p *ParallelPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, planName, p)
}
