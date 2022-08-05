package platformfee

import "github.com/jamestrandung/go-cte/cte"

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

type PlatformFee cte.SyncSideEffect

func (c PlatformFee) CTEMetadata() any {
	return struct {
		computer computer
		inout    plan
	}{}
}
