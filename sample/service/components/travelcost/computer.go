package travelcost

import (
	"context"

	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.RegisterImpureComputer(TravelCost{}, computer{})
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	return c.calculateTravelCost(casted), nil
}

func (computer) calculateTravelCost(in Input) float64 {
	return in.GetBaseCost() + in.GetTravelDuration()*in.GetCostPerMinute() + in.GetTravelDistance()*in.GetCostPerKilometer()
}
