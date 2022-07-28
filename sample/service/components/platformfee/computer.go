package platformfee

import (
	"context"
	"github.com/jamestrandung/go-cte/cte"

	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.RegisterSideEffectComputer(PlatformFee{}, computer{})
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) error {
	casted := p.(plan)

	c.addPlatformFee(casted)

	return nil
}

func (computer) addPlatformFee(p plan) {
	p.SetTotalCost(p.GetTotalCost() + p.GetPlatformFee())
}
