package travelplan

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/service/travelplan/dummy"
)

type plan interface {
	input
	output
}

type input interface {
	GetPointA() string
	GetPointB() string
}

type output interface {
	SetTravelPlan(die.AsyncResult)
}

type TravelPlan die.AsyncResult

func (p TravelPlan) GetTravelDistance() float64 {
	result := die.Outcome[dummy.TravelPlan](p.Task)
	return result.Distance
}

func (p TravelPlan) GetTravelDuration() float64 {
	result := die.Outcome[dummy.TravelPlan](p.Task)
	return result.Duration
}
