package loading

import (
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
