package costconfigs

import (
	"context"

	"github.com/jamestrandung/go-die/sample/dependencies/configsfetcher"

	"github.com/jamestrandung/go-die/sample/config"
)

func init() {
	// fmt.Println("costconfigs")
	config.Engine.RegisterImpureComputer(CostConfigs{}, computer{})
	// fmt.Println(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	return c.doFetch(casted), nil
}

func (c computer) doFetch(p plan) configsfetcher.MergedCostConfigs {
	return p.GetConfigsFetcher().Fetch()
}
