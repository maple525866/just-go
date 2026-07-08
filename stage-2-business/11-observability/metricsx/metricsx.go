package metricsx

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Registry struct {
	mu         sync.Mutex
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
}

type Counter struct {
	name string
	help string
	mu   sync.Mutex
	v    float64
}

type Gauge struct {
	name string
	help string
	mu   sync.Mutex
	v    float64
}

type Histogram struct {
	name    string
	help    string
	buckets []float64
	mu      sync.Mutex
	counts  []int
	count   int
	sum     float64
}

func NewRegistry() *Registry {
	return &Registry{counters: map[string]*Counter{}, gauges: map[string]*Gauge{}, histograms: map[string]*Histogram{}}
}

func (r *Registry) Counter(name, help string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{name: name, help: help}
	r.counters[name] = c
	return c
}

func (r *Registry) Gauge(name, help string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{name: name, help: help}
	r.gauges[name] = g
	return g
}

func (r *Registry) Histogram(name, help string, buckets []float64) *Histogram {
	r.mu.Lock()
	defer r.mu.Unlock()
	if h, ok := r.histograms[name]; ok {
		return h
	}
	sorted := append([]float64(nil), buckets...)
	sort.Float64s(sorted)
	h := &Histogram{name: name, help: help, buckets: sorted, counts: make([]int, len(sorted))}
	r.histograms[name] = h
	return h
}

func (c *Counter) Add(delta float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v += delta
}

func (g *Gauge) Set(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.v = v
}

func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.count++
	h.sum += v
	for i, b := range h.buckets {
		if v <= b {
			h.counts[i]++
		}
	}
}

func (r *Registry) Exposition() string {
	counters, gauges, histograms := r.snapshot()

	var b strings.Builder
	for _, c := range counters {
		c.mu.Lock()
		fmt.Fprintf(&b, "# HELP %s %s\n# TYPE %s counter\n%s %g\n", c.name, c.help, c.name, c.name, c.v)
		c.mu.Unlock()
	}
	for _, g := range gauges {
		g.mu.Lock()
		fmt.Fprintf(&b, "# HELP %s %s\n# TYPE %s gauge\n%s %g\n", g.name, g.help, g.name, g.name, g.v)
		g.mu.Unlock()
	}
	for _, h := range histograms {
		h.mu.Lock()
		fmt.Fprintf(&b, "# HELP %s %s\n# TYPE %s histogram\n", h.name, h.help, h.name)
		for i, upper := range h.buckets {
			fmt.Fprintf(&b, "%s_bucket{le=\"%g\"} %d\n", h.name, upper, h.counts[i])
		}
		fmt.Fprintf(&b, "%s_bucket{le=\"+Inf\"} %d\n%s_sum %g\n%s_count %d\n", h.name, h.count, h.name, h.sum, h.name, h.count)
		h.mu.Unlock()
	}
	return b.String()
}

func (r *Registry) snapshot() ([]*Counter, []*Gauge, []*Histogram) {
	r.mu.Lock()
	defer r.mu.Unlock()

	counterNames := keys(r.counters)
	gaugeNames := keys(r.gauges)
	histogramNames := keys(r.histograms)

	counters := make([]*Counter, 0, len(counterNames))
	for _, name := range counterNames {
		counters = append(counters, r.counters[name])
	}
	gauges := make([]*Gauge, 0, len(gaugeNames))
	for _, name := range gaugeNames {
		gauges = append(gauges, r.gauges[name])
	}
	histograms := make([]*Histogram, 0, len(histogramNames))
	for _, name := range histogramNames {
		histograms = append(histograms, r.histograms[name])
	}
	return counters, gauges, histograms
}

func keys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
