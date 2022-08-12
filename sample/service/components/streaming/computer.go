package streaming

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/config"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) error {
	casted := p.(inout)

	c.stream(casted)

	return nil
}

func (computer) stream(p inout) {
	config.Print("Streaming calculated cost:", p.GetTotalCost())
}
