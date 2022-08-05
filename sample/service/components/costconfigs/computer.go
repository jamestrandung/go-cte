package costconfigs

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"

	"github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"
)

type Computer struct{}

func (c Computer) Metadata() any {
    return struct {
        key   CostConfigs
        inout plan
    }{}
}

func (c Computer) Compute(ctx context.Context, p cte.MasterPlan) (any, error) {
    casted := p.(plan)

    return c.doFetch(casted), nil
}

func (c Computer) doFetch(p plan) configsfetcher.MergedCostConfigs {
    return p.GetConfigsFetcher().Fetch()
}
