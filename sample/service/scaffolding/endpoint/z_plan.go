package endpoint

import "github.com/jamestrandung/go-cte/sample/config"

var planName string

func init() {
    planName = config.Engine.AnalyzePlan(&SequentialPlan{})
}
