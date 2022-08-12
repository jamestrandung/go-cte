package streaming

import "github.com/jamestrandung/go-cte/cte"

type inout interface {
	Input
}

type Input interface {
	GetTotalCost() float64
}

type CostStreaming cte.SideEffect

func (c CostStreaming) CTEMetadata() any {
	return struct {
		computer computer
		inout    inout
	}{}
}
