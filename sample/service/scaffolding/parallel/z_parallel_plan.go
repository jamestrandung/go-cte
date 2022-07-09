package parallel

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/components/quote"
	"github.com/jamestrandung/go-die/sample/service/components/travelcost"
	"github.com/jamestrandung/go-die/sample/service/components/travelplan"
)

var planName string

func init() {
	// config.Print("ParallelPlan")
	planName = config.Engine.AnalyzePlan(&ParallelPlan{})
}

func (p *ParallelPlan) SetCostConfigs(r die.Result) {
	p.CostConfigs = (costconfigs.CostConfigs)(r)
}

func (p *ParallelPlan) SetTravelPlan(r die.Result) {
	p.TravelPlan = (travelplan.TravelPlan)(r)
}

func (p *ParallelPlan) SetTravelCost(r die.Result) {
	p.TravelCost = (travelcost.TravelCost)(r)
}

func (p *ParallelPlan) SetCalculatedCost(r die.Result) {
	p.CalculatedCost = (quote.CalculatedCost)(r)
}
