package costconfigs

import (
	"context"
	"github.com/jamestrandung/go-cte/cte"

	"github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"

	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.RegisterImpureComputer(CostConfigs{}, computer{})
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
	casted := p.(plan)

	return c.doFetch(casted), nil
}

func (c computer) doFetch(p plan) configsfetcher.MergedCostConfigs {
	return p.GetConfigsFetcher().Fetch()
}
