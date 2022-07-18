package streaming

import "github.com/jamestrandung/go-cte/cte"

type plan interface {
	Input
}

type Input interface {
	GetTotalCost() float64
	GetVATAmount() float64
}

type CostStreaming cte.SideEffectKey
