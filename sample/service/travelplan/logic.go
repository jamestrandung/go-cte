package travelplan

import "github.com/jamestrandung/go-die/sample/service/travelplan/dummy"

func (c computer) buildTravelPlan(p plan) (dummy.TravelPlan, error) {
	return c.mapService.BuildTravelPlan(p.GetPointA(), p.GetPointB())
}
