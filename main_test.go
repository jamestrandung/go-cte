package main

import (
	"context"
	"testing"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/server"
	"github.com/jamestrandung/go-cte/sample/service/miscellaneous"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/parallel"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/sequential"
)

func BenchmarkCustomPostHook_PostExecute(b *testing.B) {
	server.Serve()

	config.Engine.ConnectPostHook(&sequential.SequentialPlan{}, customPostHook{})

	p := parallel.NewPlan(
		miscellaneous.CostRequest{
			PointA: "Clementi",
			PointB: "Changi Airport",
		},
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := p.Execute(context.Background()); err != nil {
			config.Print(err)
		}
	}
}
