package metricsx

import (
	"strconv"
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

func TestRegistryExpositionConcurrentRegistration(t *testing.T) {
	reg := NewRegistry()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 200; i++ {
			reg.Counter("dynamic_counter_"+strconv.Itoa(i), "Dynamic counter").Add(1)
			reg.Gauge("dynamic_gauge_"+strconv.Itoa(i), "Dynamic gauge").Set(float64(i))
			reg.Histogram("dynamic_histogram_"+strconv.Itoa(i), "Dynamic histogram", []float64{1, 10}).Observe(float64(i))
		}
	}()
	for {
		select {
		case <-done:
			_ = reg.Exposition()
			return
		default:
			_ = reg.Exposition()
		}
	}
}
