package dummy

import (
	"github.com/jamestrandung/go-die/sample/config"
)

var planName string

func init() {
	// config.Print("DummyPlan")
	planName = config.Engine.AnalyzePlan(&DummyPlan{})
}

func (p *DummyPlan) SetSomething(something float64) {
	p.something = something
}
