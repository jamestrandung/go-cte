package loading

import (
	"context"
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte/sample/service/components/travelcost"
	"github.com/jamestrandung/go-cte/sample/service/components/travelplan"
)

type ParallelPlan struct {
	costconfigs.CostConfigs
	travelplan.TravelPlan
	travelcost.TravelCost
}

func (p *ParallelPlan) IsSequential() bool {
	return false
}

func (p *ParallelPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, planName, p)
}
