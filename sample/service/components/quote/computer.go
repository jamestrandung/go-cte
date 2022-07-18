package quote

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/fixedcost"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/sequential"
)

func init() {
	// fmt.Println("costconfigs")
	config.Engine.RegisterSwitchComputer(CalculatedCost{}, computer{})
	// fmt.Println(config.Engine)
}

type computer struct{}

// TODO: Due to pre execution can return nil, clients must take care of handling nil plan in getters
func (c computer) Switch(ctx context.Context, p any) (cte.MasterPlan, error) {
	casted := p.(plan)

	if casted.GetIsFixedCostEnabled() {
		return fixedcost.NewPlan(casted), nil
	}

	return sequential.NewPlan(casted), nil
}
