package dummy

import (
	"github.com/jamestrandung/go-die/sample/service/platformfee"
	"github.com/jamestrandung/go-die/sample/service/vat"
)

type DummyPlan struct {
	something float64
	platformfee.PlatformFee
	vat.Amount
}

func (p *DummyPlan) IsSequential() bool {
	return true
}
