package platformfee

import "github.com/jamestrandung/go-die/die"

type plan interface {
	input
	output
}

type input interface {
	GetPlatformFee() float64
	GetTotalCost() float64
}

type output interface {
	SetTotalCost(float64)
}

type PlatformFee die.ComputerKey
