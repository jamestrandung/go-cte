package streaming

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/config"
)

type Computer struct{}

func (c Computer) Metadata() any {
	return struct {
		key   CostStreaming
		inout plan
	}{}
}

func (c Computer) Compute(ctx context.Context, p cte.MasterPlan) error {
	casted := p.(plan)

	c.stream(casted)

	return nil
}

func (Computer) stream(p plan) {
	config.Print("Streaming calculated cost:", p.GetTotalCost())
}
