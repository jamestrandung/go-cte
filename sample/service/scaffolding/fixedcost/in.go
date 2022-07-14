package fixedcost

import "github.com/jamestrandung/go-cte/sample/service/components/vat"

type Input interface {
	vat.Input
	GetFixedCost() float64
}
