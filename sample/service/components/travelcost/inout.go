package travelcost

import (
	"github.com/jamestrandung/go-cte/cte"
)

type plan interface {
	input
	output
}

type input interface {
	GetBaseCost() float64
	GetTravelDistance() float64
	GetTravelDuration() float64
	GetCostPerKilometer() float64
	GetCostPerMinute() float64
}

type output interface {
	SetTravelCost(cte.Result)
}

type TravelCost cte.Result

func (r TravelCost) GetTravelCost() float64 {
	result := cte.Outcome[float64](r.Task)
	return result
}
