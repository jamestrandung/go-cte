package sequential

import (
	"github.com/jamestrandung/go-cte/sample/service/components/platformfee"
	"github.com/jamestrandung/go-cte/sample/service/components/vat"
)

type Input interface {
	preIn
	vat.Input
	platformfee.Input
}
