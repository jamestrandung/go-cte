package parallel

import (
    "context"

    "github.com/jamestrandung/go-die/sample/service/components/quote"

    "github.com/jamestrandung/go-die/sample/config"
    "github.com/jamestrandung/go-die/sample/service/components/costconfigs"
    "github.com/jamestrandung/go-die/sample/service/components/travelcost"
    "github.com/jamestrandung/go-die/sample/service/components/travelplan"
)

type ParallelPlan struct {
	Request
	Dependencies
	costconfigs.CostConfigs
	travelplan.TravelPlan
	travelcost.TravelCost
	quote.CalculatedCost
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
