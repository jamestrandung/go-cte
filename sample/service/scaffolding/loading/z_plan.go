package loading

import (
	"github.com/jamestrandung/go-cte/sample/config"
)

func init() {
	config.Engine.AnalyzePlan(&ParallelPlan{})
}
