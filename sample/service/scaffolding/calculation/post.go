package calculation

import (
    "github.com/jamestrandung/go-cte/cte"
    "github.com/jamestrandung/go-cte/sample/config"
)

type post interface {
    GetTotalCost() float64
}

type postHook struct{}

func (postHook) CTEMetadata() any {
    return struct {
        inout post
    }{}
}

func (postHook) PostExecute(p cte.Plan) error {
    config.Print("After executing sequential plan")
    casted := p.(post)

    config.Print("Calculated total cost:", casted.GetTotalCost())

    return nil
}
