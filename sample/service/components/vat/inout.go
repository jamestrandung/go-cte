package vat

import (
	"github.com/jamestrandung/go-die/die"
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
	SetVATAmount(die.SyncResult)
}

type VATAmount die.SyncResult

func (a VATAmount) GetVATAmount() float64 {
	return a.Outcome.(float64)
}
