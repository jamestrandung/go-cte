package vat

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
	casted := p.(inout)

	vatAmount := casted.GetTotalCost() * casted.GetVATPercent() / 100
	casted.SetTotalCost(casted.GetTotalCost() + vatAmount)

	return vatAmount, nil
}
