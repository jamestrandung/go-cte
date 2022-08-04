package vat

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
		key   VATAmount
		inout plan
	}{}
}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
	casted := p.(plan)

	vatAmount := casted.GetTotalCost() * casted.GetVATPercent() / 100
	casted.SetTotalCost(casted.GetTotalCost() + vatAmount)

	return vatAmount, nil
}
