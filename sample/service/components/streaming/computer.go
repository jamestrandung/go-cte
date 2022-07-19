package streaming

import (
	"context"
	"fmt"
	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.RegisterSideEffectComputer(CostStreaming{}, computer{})
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) error {
	casted := p.(plan)

	c.stream(casted)

	return nil
}

func (computer) stream(p plan) {
	fmt.Println("Streaming calculated cost:", p.GetTotalCost())
}
