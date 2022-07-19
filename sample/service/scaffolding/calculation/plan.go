package calculation

import (
	"context"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/service/components/platformfee"
	"github.com/jamestrandung/go-cte/sample/service/components/vat"
)

type SequentialPlan struct {
	Input
	preHook
	totalCost float64
	platformfee.PlatformFee
	vat.VATAmount
	postHook
	anotherPostHook
}

func NewPlan(in Input) *SequentialPlan {
	return &SequentialPlan{
		Input: in,
	}
}

func (p *SequentialPlan) IsSequential() bool {
	return true
}

func (p *SequentialPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, planName, p)
}
