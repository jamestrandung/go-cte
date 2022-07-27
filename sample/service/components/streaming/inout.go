package streaming

import "github.com/jamestrandung/go-cte/cte"

type plan interface {
	Input
}

type Input interface {
	GetTotalCost() float64
}

type CostStreaming cte.SideEffect
