package loading

import (
	"github.com/jamestrandung/go-cte/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte/sample/service/components/travelplan"
)

type Dependencies interface {
	costconfigs.Dependencies
	travelplan.Dependencies
}
