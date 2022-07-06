package platformfee

import "github.com/jamestrandung/go-die/die"

type plan interface {
	Input
	Output
}

type Input interface {
	GetPlatformFee() float64
	GetTotalCost() float64
}

type Output interface {
	SetTotalCost(float64)
}

type PlatformFee die.SideEffectKey
