package costconfigs

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/dependencies/configsfetcher"
)

type plan interface {
	input
	output
}

type Dependencies interface {
	GetConfigsFetcher() configsfetcher.Fetcher
}

type input interface {
	Dependencies
}

type output interface {
	SetCostConfigs(die.AsyncResult)
}

type CostConfigs die.AsyncResult

func (c CostConfigs) GetBaseCost() float64 {
	result := die.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.BaseCost
}

func (c CostConfigs) GetCostPerKilometer() float64 {
	result := die.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.CostPerKilometer
}

func (c CostConfigs) GetCostPerMinute() float64 {
	result := die.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.CostPerMinute
}

func (c CostConfigs) GetPlatformFee() float64 {
	result := die.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.PlatformFee
}

func (c CostConfigs) GetVATPercent() float64 {
	result := die.Outcome[configsfetcher.MergedCostConfigs](c.Task)
	return result.VATPercent
}
