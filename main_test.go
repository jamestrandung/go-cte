package main

import (
	"context"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/endpoint"
	"testing"

	"github.com/jamestrandung/go-cte/sample/dto"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/server"
)

func BenchmarkCustomPostHook_PostExecute(b *testing.B) {
    server.Serve()

    p := endpoint.NewPlan(
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
