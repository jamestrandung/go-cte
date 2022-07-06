package travelplan

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/dependencies/mapservice"
)

type plan interface {
	Input
	Output
}

type Dependencies interface {
	GetMapService() mapservice.Service
}

type Input interface {
	Dependencies
	GetPointA() string
	GetPointB() string
}

type Output interface {
	SetTravelPlan(die.AsyncResult)
}

type TravelPlan die.AsyncResult

func (p TravelPlan) GetTravelDistance() float64 {
	result := die.Outcome[mapservice.Route](p.Task)
	return result.Distance
}

func (p TravelPlan) GetTravelDuration() float64 {
	result := die.Outcome[mapservice.Route](p.Task)
	return result.Duration
}
