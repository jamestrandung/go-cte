package sequential

import (
	"github.com/jamestrandung/go-die/sample/service/components/platformfee"
	"github.com/jamestrandung/go-die/sample/service/components/vat"
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
