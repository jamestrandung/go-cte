package platformfee

import (
	"context"

	"github.com/jamestrandung/go-die/sample/config"
)

// Computers without any external dependencies can register itself directly
// with the engine using init()
func init() {
	// config.Print("platformfee")
	config.Engine.RegisterSideEffectComputer(PlatformFee{}, computer{})
	// config.Print(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) error {
	casted := p.(plan)

	c.addPlatformFee(casted)

	return nil
}

func (computer) addPlatformFee(p plan) {
	p.SetTotalCost(p.GetTotalCost() + p.GetPlatformFee())
}
