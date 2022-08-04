package streaming

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
		key   CostStreaming
		inout plan
	}{}
}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) error {
	casted := p.(plan)

	c.stream(casted)

	return nil
}

func (computer) stream(p plan) {
	config.Print("Streaming calculated cost:", p.GetTotalCost())
}
