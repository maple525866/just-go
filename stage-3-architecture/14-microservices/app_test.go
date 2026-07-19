package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRunCompletesProductDetailsFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var out bytes.Buffer
	if err := run(ctx, &out); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"book-1", "Go Microservices", `"quantity":10`, `"stock_version":1`} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q missing %q", got, want)
		}
	}
}

func TestRunHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := run(ctx, &bytes.Buffer{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestRunReportsServeFailure(t *testing.T) {
	want := errors.New("accept failed")
	listen := func(network, address string) (net.Listener, error) {
		listener, err := net.Listen(network, address)
		if err != nil {
			return nil, err
		}
		return &failingListener{Listener: listener, err: want}, nil
	}

	err := runWithListen(context.Background(), &bytes.Buffer{}, listen)
	if err == nil || !strings.Contains(err.Error(), "serve") || !errors.Is(err, want) {
		t.Fatalf("error = %v, want wrapped serve failure", err)
	}
}

func TestRunStopsAfterPostStartCancellation(t *testing.T) {
	started := make(chan struct{})
	listenCount := 0
	listen := func(network, address string) (net.Listener, error) {
		listener, err := net.Listen(network, address)
		if err != nil {
			return nil, err
		}
		listenCount++
		if listenCount == 3 {
			close(started)
			return &gatedListener{Listener: listener, gate: make(chan struct{})}, nil
		}
		return listener, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() { result <- runWithListen(ctx, &bytes.Buffer{}, listen) }()
	<-started
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("run did not stop after cancellation")
	}
}

type failingListener struct {
	net.Listener
	err error
}

func (l *failingListener) Accept() (net.Conn, error) {
	return nil, l.err
}

type gatedListener struct {
	net.Listener
	gate chan struct{}
	once sync.Once
}

func (l *gatedListener) Accept() (net.Conn, error) {
	<-l.gate
	return nil, net.ErrClosed
}

func (l *gatedListener) Close() error {
	l.once.Do(func() { close(l.gate) })
	return l.Listener.Close()
}
