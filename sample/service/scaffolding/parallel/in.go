package parallel

import (
	"github.com/jamestrandung/go-die/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/components/travelplan"
)

type Dependencies interface {
	costconfigs.Dependencies
	travelplan.Dependencies
}

type Request interface {
	GetPointA() string
	GetPointB() string
}
