package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"just-go/stage-2-business/11-observability/healthx"
	"just-go/stage-2-business/11-observability/loggingx"
	"just-go/stage-2-business/11-observability/metricsx"
	"just-go/stage-2-business/11-observability/tracex"
)

type Router struct {
	checker *healthx.Checker
	metrics *metricsx.Registry
	logger  *slog.Logger
}

func NewRouter(checker *healthx.Checker, reg *metricsx.Registry) http.Handler {
	return &Router{checker: checker, metrics: reg, logger: loggingx.NewJSONLogger(os.Stdout, slog.LevelInfo)}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	r.metrics.Counter("http_requests_total", "Total HTTP requests").Add(1)
	defer func() {
		r.metrics.Histogram("http_request_duration_seconds", "HTTP request duration", []float64{0.01, 0.05, 0.1, 0.5, 1}).Observe(time.Since(start).Seconds())
	}()

	traceID := req.Header.Get("X-Trace-ID")
	ctx := req.Context()
	if traceID != "" {
		ctx = tracex.WithTraceID(ctx, traceID)
	}
	ctx, span := tracex.StartSpan(ctx, req.Method+" "+req.URL.Path)
	w.Header().Set("X-Trace-ID", span.TraceID)
	loggingx.Info(ctx, r.logger, "http.request", slog.String("path", req.URL.Path), slog.String("span_id", span.SpanID))

	switch req.URL.Path {
	case "/livez":
		writeReport(w, r.checker.Liveness(ctx))
	case "/readyz":
		report := r.checker.Readiness(ctx)
		if !report.OK {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		writeJSON(w, report)
	case "/metrics":
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(r.metrics.Exposition()))
	case "/work":
		r.metrics.Counter("work_units_total", "Total demo work units").Add(1)
		writeJSON(w, map[string]string{"status": "done", "trace_id": span.TraceID, "span_id": span.SpanID})
	default:
		http.NotFound(w, req)
	}
}

func writeReport(w http.ResponseWriter, report healthx.Report) {
	writeJSON(w, report)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
