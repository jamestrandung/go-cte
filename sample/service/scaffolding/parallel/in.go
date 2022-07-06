package parallel

import (
	"github.com/jamestrandung/go-die/sample/service/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/travelplan"
)

type Dependencies interface {
	costconfigs.Dependencies
	travelplan.Dependencies
}

type Request interface {
	GetPointA() string
	GetPointB() string
}
