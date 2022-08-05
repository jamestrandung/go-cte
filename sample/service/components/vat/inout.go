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
}

type VATAmount cte.SyncResult

func (a VATAmount) CTEMetadata() any {
	return struct {
		computer computer
		inout    plan
	}{}
}

func (a VATAmount) GetVATAmount() float64 {
	return cte.Cast[float64](a.Outcome)
}
