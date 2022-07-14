package main

import (
	"context"
	"testing"

	"github.com/jamestrandung/go-cte/sample/dto"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/server"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/parallel"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/sequential"
)

func BenchmarkCustomPostHook_PostExecute(b *testing.B) {
	server.Serve()

	config.Engine.ConnectPostHook(&sequential.SequentialPlan{}, customPostHook{})

	p := parallel.NewPlan(
		dto.CostRequest{
			PointA: "Clementi",
			PointB: "Changi Airport",
		},
		server.Dependencies,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := p.Execute(context.Background()); err != nil {
			config.Print(err)
		}
	}
}
