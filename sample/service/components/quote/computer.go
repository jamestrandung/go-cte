package quote

import (
	"context"

	"github.com/jamestrandung/go-cte/cte"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/calculation"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/fixedcost"
)

type Computer struct{}

func (c Computer) Metadata() any {
    return struct {
        key   FixedCostBranch
        inout plan
    }{}
}

// TODO: Due to pre execution can return nil, clients must take care of handling nil plan in getters
func (c Computer) Switch(ctx context.Context, p cte.MasterPlan) (cte.MasterPlan, error) {
    casted := p.(plan)

    if casted.GetIsFixedCostEnabled() {
        return fixedcost.NewPlan(casted), nil
    }

    return calculation.NewPlan(casted), nil
}
