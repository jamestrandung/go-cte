package parallel

import (
	"github.com/jamestrandung/go-cte/sample/config"
)

var planName string

func init() {
	// config.Print("ParallelPlan")
	planName = config.Engine.AnalyzePlan(&ParallelPlan{})
}
