package main

import (
	"context"
	"fmt"

	"just-go/stage-2-business/11-observability/healthx"
	"just-go/stage-2-business/11-observability/metricsx"
	"just-go/stage-2-business/11-observability/tracex"
)

func main() {
	ctx, span := tracex.StartSpan(context.Background(), "chapter-11.demo")
	reg := metricsx.NewRegistry()
	reg.Counter("demo_runs_total", "Demo runs").Add(1)
	reg.Gauge("demo_readiness", "Demo readiness state").Set(1)
	checker := healthx.NewChecker()
	checker.AddLiveness("process", healthx.OK)
	checker.AddReadiness("dependencies", healthx.OK)

	fmt.Println("Chapter 11 Observability")
	fmt.Printf("trace: %s span: %s (%s)\n", span.TraceID, span.SpanID, span.Name)
	fmt.Printf("current trace from context: %s\n", tracex.TraceID(ctx))
	fmt.Printf("liveness ok: %t readiness ok: %t\n", checker.Liveness(ctx).OK, checker.Readiness(ctx).OK)
	fmt.Println("metrics:")
	fmt.Print(reg.Exposition())
}
