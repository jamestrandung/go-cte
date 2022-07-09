package travelcost

import (
	"github.com/jamestrandung/go-die/die"
)

type plan interface {
	Input
}

type Input interface {
	GetBaseCost() float64
	GetTravelDistance() float64
	GetTravelDuration() float64
	GetCostPerKilometer() float64
	GetCostPerMinute() float64
}

type TravelCost die.Result

func (r TravelCost) GetTravelCost() float64 {
	result := die.Outcome[float64](r.Task)
	return result
}
