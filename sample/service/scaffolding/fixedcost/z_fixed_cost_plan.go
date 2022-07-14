package fixedcost

import (
	"github.com/jamestrandung/go-cte/sample/config"
)

var planName string

func init() {
	// config.Print("SequentialPlan")
	planName = config.Engine.AnalyzePlan(&FixedCostPlan{})
}

func (p *FixedCostPlan) GetTotalCost() float64 {
	return p.totalCost
}

func (p *FixedCostPlan) SetTotalCost(totalCost float64) {
	p.totalCost = totalCost
}
