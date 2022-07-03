package travelcost

import (
	"github.com/jamestrandung/go-die/die"
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
	SetTravelCost(die.AsyncResult)
}

type TravelCost die.AsyncResult

func (r TravelCost) GetTravelCost() float64 {
	result := die.Outcome[float64](r.Task)
	return result
}
