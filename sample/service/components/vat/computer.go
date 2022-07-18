package vat

import (
	"context"

	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	// config.Print("vat")
	config.Engine.RegisterImpureComputer(VATAmount{}, computer{})
	// config.Print(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	vatAmount := casted.GetTotalCost() * casted.GetVATPercent() / 100
	casted.SetTotalCost(casted.GetTotalCost() + vatAmount)

	return vatAmount, nil
}
