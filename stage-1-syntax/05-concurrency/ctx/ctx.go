// Package ctx 演示 context 的取消、超时与协作式退出。
//
// import 路径：just-go/stage-1-syntax/05-concurrency/ctx
package ctx

import (
	"context"
	"errors"
	"time"
)

// CancellationDemo 启动一个观察 ctx.Done 的 worker，并通过 cancel 让它退出。
func CancellationDemo() string {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan string, 1)
	go func() {
		<-ctx.Done()
		done <- "cancelled"
	}()
	cancel()
	return <-done
}

// TimeoutDemo 等待 context timeout，并返回是否得到 deadline exceeded。
func TimeoutDemo(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	<-ctx.Done()
	return errors.Is(ctx.Err(), context.DeadlineExceeded)
}

// WorkerWithContext 模拟协作式 worker：收到任务则处理，收到取消信号则退出。
func WorkerWithContext(parent context.Context, jobs <-chan string) string {
	select {
	case <-parent.Done():
		return parent.Err().Error()
	case job := <-jobs:
		return "processed:" + job
	}
}
