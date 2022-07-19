package vat

import (
	"context"

	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.RegisterImpureComputer(VATAmount{}, computer{})
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	vatAmount := casted.GetTotalCost() * casted.GetVATPercent() / 100
	casted.SetTotalCost(casted.GetTotalCost() + vatAmount)

	return vatAmount, nil
}
