package travelcost

import (
	"github.com/jamestrandung/go-cte/cte"
)

type inout interface {
	Input
}

type Input interface {
	GetBaseCost() float64
	GetTravelDistance() float64
	GetTravelDuration() float64
	GetCostPerKilometer() float64
	GetCostPerMinute() float64
}

type TravelCost cte.Result

func (r TravelCost) CTEMetadata() any {
	return struct {
		computer Computer
		inout    inout
	}{}
}

func (r TravelCost) GetTravelCost() float64 {
	result := cte.Outcome[float64](r.Task)
	return result
}
