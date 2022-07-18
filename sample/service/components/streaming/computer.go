package streaming

import (
	"context"
	"fmt"
	"github.com/jamestrandung/go-cte/sample/config"
)

// Computers without any external dependencies can register itself directly
// with the engine using init()
func init() {
	// config.Print("streaming")
	config.Engine.RegisterSideEffectComputer(CostStreaming{}, computer{})
	// config.Print(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) error {
	casted := p.(plan)

	c.stream(casted)

	return nil
}

func (computer) stream(p plan) {
	fmt.Println("Streaming calculated cost:", p.GetVATAmount())
}
