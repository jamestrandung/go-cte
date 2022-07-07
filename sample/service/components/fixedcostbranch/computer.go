package fixedcostbranch

import (
	"context"
	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/scaffolding/fixedcost"
	"github.com/jamestrandung/go-die/sample/service/scaffolding/sequential"
)

func init() {
	// fmt.Println("costconfigs")
	config.Engine.RegisterImpureComputer(CalculatedCost{}, computer{})
	// fmt.Println(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	if casted.GetIsFixedCostEnabled() {
		fixedCostPlan := fixedcost.NewPlan(casted)
		if err := fixedCostPlan.Execute(ctx); err != nil {
			return nil, err
		}

		return (result)(fixedCostPlan), nil
	}

	sequentialPlan := sequential.NewPlan(casted)
	if err := sequentialPlan.Execute(ctx); err != nil {
		return nil, err
	}

	return (result)(sequentialPlan), nil
}
