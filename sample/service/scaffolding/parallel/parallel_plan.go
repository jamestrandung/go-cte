package parallel

import (
	"context"

	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/miscellaneous"
	"github.com/jamestrandung/go-die/sample/service/scaffolding/sequential"
	"github.com/jamestrandung/go-die/sample/service/travelcost"
	"github.com/jamestrandung/go-die/sample/service/travelplan"
)

type ParallelPlan struct {
	miscellaneous.CostRequest
	costconfigs.CostConfigs
	travelplan.TravelPlan
	travelcost.TravelCost
	sequential.SequentialPlan
}

func NewPlan(r miscellaneous.CostRequest) *ParallelPlan {
	return &ParallelPlan{
		CostRequest: r,
	}
}

func (p *ParallelPlan) IsSequential() bool {
	return false
}

func (p *ParallelPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, planName, p)
}
