package vat

import (
	"github.com/jamestrandung/go-cte/cte"
)

type inout interface {
	Input
	Output
}

type dummy interface {
	GetVATPercent() float64
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
		inout    inout
	}{}
}

func (a VATAmount) GetVATAmount() float64 {
	return cte.Cast[float64](a.Outcome)
}
