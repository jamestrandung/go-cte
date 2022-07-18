package travelplan

import (
	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/dependencies/mapservice"
)

type plan interface {
	Input
}

type Dependencies interface {
	GetMapService() mapservice.Service
}

type Input interface {
	Dependencies
	GetPointA() string
	GetPointB() string
}

type TravelPlan cte.Result

func (p TravelPlan) GetTravelDistance() float64 {
	result := cte.Outcome[mapservice.Route](p.Task)
	return result.Distance
}

func (p TravelPlan) GetTravelDuration() float64 {
	result := cte.Outcome[mapservice.Route](p.Task)
	return result.Duration
}
