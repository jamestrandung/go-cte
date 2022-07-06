package parallel

import (
	"context"

	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/scaffolding/sequential"
	"github.com/jamestrandung/go-die/sample/service/travelcost"
	"github.com/jamestrandung/go-die/sample/service/travelplan"
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
