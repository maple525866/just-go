// Package goroutine 演示 goroutine 的启动、等待与结果收集。
//
// import 路径：just-go/stage-1-syntax/05-concurrency/goroutine
package goroutine

import "sync"

// Result 表示一个并发任务的结果。
type Result struct {
	ID    int
	Value int
}

// RunSquares 为每个输入启动一个 goroutine 计算平方，并等待全部任务完成后返回结果。
func RunSquares(inputs []int) []Result {
	results := make([]Result, len(inputs))
	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for i, n := range inputs {
		i, n := i, n
		go func() {
			defer wg.Done()
			results[i] = Result{ID: i, Value: n * n}
		}()
	}
	wg.Wait()
	return results
}

// LifecycleSummary 返回 goroutine 生命周期的关键步骤摘要。
func LifecycleSummary() []string {
	return []string{"start goroutine", "do work", "signal done", "wait and collect"}
}
