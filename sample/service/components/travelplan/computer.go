package travelplan

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"

	"github.com/jamestrandung/go-cte/sample/dependencies/mapservice"

	"github.com/jamestrandung/go-cte/sample/config"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
	casted := p.(plan)

	route, err := casted.GetMapService().GetRoute(casted.GetPointA(), casted.GetPointB())
	if err != nil {
		return c.calculateStraightLineDistance(casted), nil
	}

	return route, nil
}

func (c computer) calculateStraightLineDistance(p plan) mapservice.Route {
	config.Printf("Building route from %s to %s using straight-line distance\n", p.GetPointA(), p.GetPointB())
	return mapservice.Route{
		Distance: 4,
		Duration: 5,
	}
}
