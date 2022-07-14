package vat

import (
	"github.com/jamestrandung/go-cte/cte"
)

type plan interface {
	Input
	Output
}

type Input interface {
	GetVATPercent() float64
	GetTotalCost() float64
}

type Output interface {
	SetTotalCost(float64)
	SetVATAmount(cte.SyncResult)
}

type VATAmount cte.SyncResult

func (a VATAmount) GetVATAmount() float64 {
	return a.Outcome.(float64)
}
