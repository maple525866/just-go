package observability

import (
	"just-go/stage-2-business/11-observability/healthx"
	"just-go/stage-2-business/11-observability/metricsx"
)

type Bundle struct {
	Metrics *metricsx.Registry
	Health  *healthx.Checker
}

func New() Bundle {
	h := healthx.NewChecker()
	h.AddLiveness("process", healthx.OK)
	h.AddReadiness("store", healthx.OK)
	return Bundle{Metrics: metricsx.NewRegistry(), Health: h}
}
