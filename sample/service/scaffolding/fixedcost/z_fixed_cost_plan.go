package fixedcost

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/components/vat"
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

func (p *FixedCostPlan) SetVATAmount(r die.SyncResult) {
	p.VATAmount = (vat.VATAmount)(r)
}
