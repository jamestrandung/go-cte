package fixedcost

import "github.com/jamestrandung/go-die/sample/service/components/vat"

type Input interface {
	vat.Input
	GetFixedCost() float64
}
