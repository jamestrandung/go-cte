package fixedcostbranch

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/service/scaffolding/fixedcost"
	"github.com/jamestrandung/go-die/sample/service/scaffolding/sequential"
)

type plan interface {
	Input
	Output
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

type Output interface {
	SetCalculatedCost(die.AsyncResult)
}

type CalculatedCost die.AsyncResult

func (c CalculatedCost) GetTotalCost() float64 {
	return die.Outcome[result](c.Task).GetTotalCost()
}

func (c CalculatedCost) GetVATAmount() float64 {
	return die.Outcome[result](c.Task).GetVATAmount()
}
