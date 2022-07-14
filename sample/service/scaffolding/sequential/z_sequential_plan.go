package sequential

import (
	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/components/vat"
)

var planName string

func init() {
	// config.Print("SequentialPlan")
	planName = config.Engine.AnalyzePlan(&SequentialPlan{})
}

func (p *SequentialPlan) GetTotalCost() float64 {
	return p.totalCost
}

func (p *SequentialPlan) SetTotalCost(totalCost float64) {
	p.totalCost = totalCost
}

func (p *SequentialPlan) SetVATAmount(r cte.SyncResult) {
	p.VATAmount = (vat.VATAmount)(r)
}
