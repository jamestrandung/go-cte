package vat

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
)

type Computer struct{}

func (c Computer) Metadata() any {
    return struct {
        key   VATAmount
        inout plan
    }{}
}

func (c Computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
    casted := p.(plan)

    vatAmount := casted.GetTotalCost() * casted.GetVATPercent() / 100
    casted.SetTotalCost(casted.GetTotalCost() + vatAmount)

    return vatAmount, nil
}
