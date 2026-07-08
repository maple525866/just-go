package healthx

import (
	"context"
	"errors"
)

type Check func(context.Context) error

type namedCheck struct {
	name string
	fn   Check
}

type Checker struct {
	liveness  []namedCheck
	readiness []namedCheck
}

type Result struct {
	Name  string `json:"name"`
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Report struct {
	OK     bool     `json:"ok"`
	Checks []Result `json:"checks"`
}

func OK(context.Context) error { return nil }

func ErrNotReady(message string) error { return errors.New(message) }

func NewChecker() *Checker { return &Checker{} }

func (c *Checker) AddLiveness(name string, fn Check) {
	c.liveness = append(c.liveness, namedCheck{name: name, fn: fn})
}

func (c *Checker) AddReadiness(name string, fn Check) {
	c.readiness = append(c.readiness, namedCheck{name: name, fn: fn})
}

func (c *Checker) Liveness(ctx context.Context) Report { return runChecks(ctx, c.liveness) }

func (c *Checker) Readiness(ctx context.Context) Report { return runChecks(ctx, c.readiness) }

func runChecks(ctx context.Context, checks []namedCheck) Report {
	report := Report{OK: true, Checks: make([]Result, 0, len(checks))}
	for _, check := range checks {
		result := Result{Name: check.name, OK: true}
		if err := check.fn(ctx); err != nil {
			result.OK = false
			result.Error = err.Error()
			report.OK = false
		}
		report.Checks = append(report.Checks, result)
	}
	return report
}
