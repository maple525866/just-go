package metricsx

import (
	"strings"
	"testing"
)

func TestRegistryRecordsAndExposesPrometheusText(t *testing.T) {
	reg := NewRegistry()
	reg.Counter("http_requests_total", "Total HTTP requests").Add(3)
	reg.Gauge("inflight_requests", "In-flight requests").Set(2)
	reg.Histogram("request_duration_seconds", "Request duration", []float64{0.1, 0.5, 1}).Observe(0.2)
	reg.Histogram("request_duration_seconds", "Request duration", []float64{0.1, 0.5, 1}).Observe(2)

	out := reg.Exposition()
	checks := []string{
		"# TYPE http_requests_total counter",
		"http_requests_total 3",
		"# TYPE inflight_requests gauge",
		"inflight_requests 2",
		`request_duration_seconds_bucket{le="0.5"} 1`,
		`request_duration_seconds_bucket{le="+Inf"} 2`,
		"request_duration_seconds_count 2",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("exposition missing %q:\n%s", want, out)
		}
	}
}
