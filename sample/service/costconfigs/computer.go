package costconfigs

import (
	"context"

	"github.com/jamestrandung/go-die/sample/config"
	"github.com/jamestrandung/go-die/sample/service/costconfigs/dummy"
)

// Computers with external dependencies still has to register itself with the
// engine using init() so that we can perform validations on plans
func init() {
	// fmt.Println("costconfigs")
	config.Engine.RegisterImpureComputer(CostConfigs{}, computer{})
	// fmt.Println(config.Engine)
}

type computer struct {
	fetcher dummy.CostConfigsFetcher
}

// Computers with external dependencies can register itself with the engine
// via an exported InitComputer() that takes in dependencies as arguments
// to overwrite the dummy computer registered via init()

// InitComputer ...
func InitComputer(fetcher dummy.CostConfigsFetcher) {
	c := computer{
		fetcher: fetcher,
	}

	// fmt.Println("costconfigs")
	config.Engine.RegisterImpureComputer(CostConfigs{}, c)
	// fmt.Println(config.Engine)
}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	return c.doFetch(), nil
}
