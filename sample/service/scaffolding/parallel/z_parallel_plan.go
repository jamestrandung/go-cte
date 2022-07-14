package parallel

import (
	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte/sample/service/components/travelcost"
	"github.com/jamestrandung/go-cte/sample/service/components/travelplan"
)

var planName string

func init() {
	// config.Print("ParallelPlan")
	planName = config.Engine.AnalyzePlan(&ParallelPlan{})
}

func (p *ParallelPlan) SetCostConfigs(r cte.Result) {
	p.CostConfigs = (costconfigs.CostConfigs)(r)
}

func (p *ParallelPlan) SetTravelPlan(r cte.Result) {
	p.TravelPlan = (travelplan.TravelPlan)(r)
}

func (p *ParallelPlan) SetTravelCost(r cte.Result) {
	p.TravelCost = (travelcost.TravelCost)(r)
}
