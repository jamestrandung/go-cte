package costconfigs

import (
	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"
)

type plan interface {
	Dependencies
}

type Dependencies interface {
	GetConfigsFetcher() configsfetcher.Fetcher
}

type CostConfigs cte.Result

func (c CostConfigs) GetBaseCost() float64 {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.BaseCost
}

func (c CostConfigs) GetCostPerKilometer() float64 {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.CostPerKilometer
}

func (c CostConfigs) GetCostPerMinute() float64 {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.CostPerMinute
}

func (c CostConfigs) GetPlatformFee() float64 {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.PlatformFee
}

func (c CostConfigs) GetVATPercent() float64 {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.VATPercent
}

func (c CostConfigs) GetIsFixedCostEnabled() bool {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.IsFixedCostEnabled
}

func (c CostConfigs) GetFixedCost() float64 {
	result := cte.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.FixedCost
}
