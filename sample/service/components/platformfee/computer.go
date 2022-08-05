package platformfee

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
)

type Computer struct{}

func (c Computer) Metadata() any {
    return struct {
        key   PlatformFee
        inout plan
    }{}
}

func (c Computer) Compute(ctx context.Context, p cte.MasterPlan) error {
    casted := p.(plan)

    c.addPlatformFee(casted)

    return nil
}

func (Computer) addPlatformFee(p plan) {
    p.SetTotalCost(p.GetTotalCost() + p.GetPlatformFee())
}
