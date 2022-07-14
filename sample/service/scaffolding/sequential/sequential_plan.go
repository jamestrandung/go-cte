package sequential

import (
	"github.com/jamestrandung/go-cte/sample/service/components/platformfee"
	"github.com/jamestrandung/go-cte/sample/service/components/vat"
)

type SequentialPlan struct {
	preHook
	totalCost float64
	platformfee.PlatformFee
	vat.VATAmount
	postHook
	anotherPostHook
}

func (p *SequentialPlan) IsSequential() bool {
	return true
}
