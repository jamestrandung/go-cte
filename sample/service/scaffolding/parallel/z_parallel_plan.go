package parallel

import (
	"github.com/jamestrandung/go-die/sample/config"
)

var planName string

func init() {
	// config.Print("ParallelPlan")
	planName = config.Engine.AnalyzePlan(&ParallelPlan{})
}
