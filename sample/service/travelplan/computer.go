package travelplan

import (
	"context"

	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/travelplan/dummy"
)

// Computers with external dependencies still has to register itself with the
// engine using init() so that we can perform validations on plans
func init() {
	// config.Print("travelplan")
	config.Engine.RegisterImpureComputer(TravelPlan{}, computer{})
	// config.Print(config.Engine)
}

type computer struct {
	mapService dummy.MapService
}

// Computers with external dependencies can register itself with the engine
// via an exported InitComputer() that takes in dependencies as arguments
// to overwrite the dummy computer registered via init()

// InitComputer ...
func InitComputer(mapService dummy.MapService) {
	c := computer{
		mapService: mapService,
	}

	// config.Print("travelplan")
	config.Engine.RegisterImpureComputer(TravelPlan{}, c)
	// config.Print(config.Engine)
}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	travelPlan, err := c.mapService.BuildTravelPlan(casted.GetPointA(), casted.GetPointB())
	if err != nil {
		return c.calculateStraightLineDistance(casted), nil
	}

	return travelPlan, nil
}

func (c computer) calculateStraightLineDistance(p plan) dummy.TravelPlan {
	config.Printf("Building travel plan from %s to %s using straight-line distance\n", p.GetPointA(), p.GetPointB())
	return dummy.TravelPlan{
		Distance: 4,
		Duration: 5,
	}
}
