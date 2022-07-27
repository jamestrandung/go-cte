package quote

import (
	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/calculation"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/fixedcost"
)

type plan interface {
	Input
}

type Input interface {
	fixedcost.Input
	calculation.Input
	GetIsFixedCostEnabled() bool
}

type result interface {
	GetTotalCost() float64
	GetVATAmount() float64
}

type FixedCostBranch cte.Result

func (c FixedCostBranch) GetTotalCost() float64 {
	return cte.Outcome[result](c.Task).GetTotalCost()
}

func (c FixedCostBranch) GetVATAmount() float64 {
	return cte.Outcome[result](c.Task).GetVATAmount()
}
