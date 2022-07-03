package vat

import (
	"github.com/jamestrandung/go-die/die"
)

type plan interface {
	input
	output
}

type input interface {
	GetVATPercent() float64
	GetTotalCost() float64
}

type output interface {
	SetTotalCost(float64)
	SetVATAmount(die.SyncResult)
}

type VATAmount die.SyncResult

func (a VATAmount) GetVATAmount() float64 {
	return a.Outcome.(float64)
}
