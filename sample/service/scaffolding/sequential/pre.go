package sequential

import "github.com/jamestrandung/go-cte/sample/config"

type pre interface {
	preIn
	preOut
}

type preIn interface {
	GetTravelCost() float64
}

type preOut interface {
	SetTotalCost(float64)
}

type preHook struct{}

func (preHook) PreExecute(p any) error {
	config.Print("Before executing sequential plan")
	casted := p.(pre)

	casted.SetTotalCost(casted.GetTravelCost())

	return nil
}
