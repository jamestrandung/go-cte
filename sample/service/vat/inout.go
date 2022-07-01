package vat

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
	SetVATAmount(VATAmount)
}

type VATAmount struct {
	vatAmount float64
}
