package travelcost

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"

	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.RegisterComputer(computer{})
}

type computer struct{}

func (c computer) Metadata() any {
	return struct {
		key   TravelCost
		inout plan
	}{}
}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
	casted := p.(plan)

	return c.calculateTravelCost(casted), nil
}

func (computer) calculateTravelCost(in Input) float64 {
	return in.GetBaseCost() + in.GetTravelDuration()*in.GetCostPerMinute() + in.GetTravelDistance()*in.GetCostPerKilometer()
}
