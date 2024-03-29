package costconfigs

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"

	"github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
	casted := p.(inout)

	return c.doFetch(casted), nil
}

func (c computer) doFetch(p inout) configsfetcher.MergedCostConfigs {
	return p.GetConfigsFetcher().Fetch()
}
