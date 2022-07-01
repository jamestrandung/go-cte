package costconfigs

import "github.com/jamestrandung/go-die/sample/service/costconfigs/dummy"

func (c computer) doFetch() dummy.MergedCostConfigs {
	return c.fetcher.Fetch()
}
