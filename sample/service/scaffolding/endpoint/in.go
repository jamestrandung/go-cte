package endpoint

import (
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/loading"
)

type Dependencies interface {
	loading.Dependencies
}

type Request interface {
	GetPointA() string
	GetPointB() string
}
