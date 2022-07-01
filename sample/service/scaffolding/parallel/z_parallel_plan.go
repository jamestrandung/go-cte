package parallel

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/travelcost"
	"github.com/jamestrandung/go-die/sample/service/travelplan"
)

var planName string

func init() {
	// config.Print("ParallelPlan")
	planName = config.Engine.AnalyzePlan(&ParallelPlan{})
}

func (p *ParallelPlan) SetCostConfigs(o die.AsyncResult) {
	p.CostConfigs = (costconfigs.CostConfigs)(o)
}

func (p *ParallelPlan) SetTravelPlan(o die.AsyncResult) {
	p.TravelPlan = (travelplan.TravelPlan)(o)
}

func (p *ParallelPlan) SetTravelCost(o die.AsyncResult) {
	p.TravelCost = (travelcost.TravelCost)(o)
}
