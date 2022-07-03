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

func (p *ParallelPlan) SetCostConfigs(r die.AsyncResult) {
	p.CostConfigs = (costconfigs.CostConfigs)(r)
}

func (p *ParallelPlan) SetTravelPlan(r die.AsyncResult) {
	p.TravelPlan = (travelplan.TravelPlan)(r)
}

func (p *ParallelPlan) SetTravelCost(r die.AsyncResult) {
	p.TravelCost = (travelcost.TravelCost)(r)
}
