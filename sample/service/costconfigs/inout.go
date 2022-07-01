package costconfigs

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/service/costconfigs/dummy"
)

type plan interface {
	output
}

type output interface {
	SetCostConfigs(die.AsyncResult)
}

type CostConfigs die.AsyncResult

func (r CostConfigs) GetBaseCost() float64 {
	result := die.ExtractAsync[dummy.MergedCostConfigs](r.Task)
	return result.BaseCost
}

func (r CostConfigs) GetCostPerKilometer() float64 {
	result := die.ExtractAsync[dummy.MergedCostConfigs](r.Task)
	return result.CostPerKilometer
}

func (r CostConfigs) GetCostPerMinute() float64 {
	result := die.ExtractAsync[dummy.MergedCostConfigs](r.Task)
	return result.CostPerMinute
}

func (r CostConfigs) GetPlatformFee() float64 {
	result := die.ExtractAsync[dummy.MergedCostConfigs](r.Task)
	return result.PlatformFee
}

func (r CostConfigs) GetVATPercent() float64 {
	result := die.ExtractAsync[dummy.MergedCostConfigs](r.Task)
	return result.VATPercent
}
