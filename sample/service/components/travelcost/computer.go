package travelcost

import (
	"context"

	"github.com/jamestrandung/go-cte/sample/config"
)

// Computers without any external dependencies can register itself directly
// with the engine using init()
func init() {
	// config.Print("travelcost")
	config.Engine.RegisterImpureComputer(TravelCost{}, computer{})
	// config.Print(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	return c.calculateTravelCost(casted), nil
}

func (computer) calculateTravelCost(in input) float64 {
	return in.GetBaseCost() + in.GetTravelDuration()*in.GetCostPerMinute() + in.GetTravelDistance()*in.GetCostPerKilometer()
}
