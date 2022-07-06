package fixedcost

import (
    "context"
    "github.com/jamestrandung/go-die/sample/config"
    "github.com/jamestrandung/go-die/sample/service/vat"
)

type FixedCostPlan struct {
    Input
    totalCost float64
    vat.VATAmount
}

func (p *FixedCostPlan) IsSequential() bool {
    return true
}

func (p *FixedCostPlan) Execute(ctx context.Context) error {
    p.preExecute()

    return config.Engine.ExecuteMasterPlan(ctx, planName, p)
}

func (p *FixedCostPlan) preExecute() {
    p.totalCost = p.GetFixedCost()
}