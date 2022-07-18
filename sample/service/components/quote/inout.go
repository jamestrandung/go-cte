package quote

import (
	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/fixedcost"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/sequential"
)

type plan interface {
	Input
}

type Input interface {
	fixedcost.Input
	sequential.Input
	GetIsFixedCostEnabled() bool
}

type result interface {
	GetTotalCost() float64
	GetVATAmount() float64
}

type CalculatedCost cte.Result

func (c CalculatedCost) GetTotalCost() float64 {
	return cte.Outcome[result](c.Task).GetTotalCost()
}

func (c CalculatedCost) GetVATAmount() float64 {
	return cte.Outcome[result](c.Task).GetVATAmount()
}
