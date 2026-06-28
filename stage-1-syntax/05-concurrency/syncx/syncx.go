// Package syncx 演示 sync 包常见原语：Mutex、RWMutex、WaitGroup、Once。
//
// import 路径：just-go/stage-1-syntax/05-concurrency/syncx
package syncx

import "sync"

// SafeCounter 用 Mutex 保护共享计数。
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// Add 安全增加计数。
func (c *SafeCounter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

// Value 安全读取计数。
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// CountWithMutex 启动 n 个 goroutine 并发累加，使用 Mutex 避免 data race。
func CountWithMutex(n int) int {
	var counter SafeCounter
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	return counter.Value()
}

// ScoreBoard 用 RWMutex 演示多读少写的共享状态保护。
type ScoreBoard struct {
	mu     sync.RWMutex
	scores map[string]int
}

// NewScoreBoard 创建分数表。
func NewScoreBoard() *ScoreBoard {
	return &ScoreBoard{scores: make(map[string]int)}
}

// Set 写入分数。
func (b *ScoreBoard) Set(name string, score int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.scores[name] = score
}

// Get 读取分数。
func (b *ScoreBoard) Get(name string) (int, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	score, ok := b.scores[name]
	return score, ok
}

// InitOnceConcurrently 多次并发触发初始化，但 sync.Once 保证初始化函数只执行一次。
func InitOnceConcurrently(times int) int {
	var once sync.Once
	var initialized int
	var wg sync.WaitGroup
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func() {
			defer wg.Done()
			once.Do(func() { initialized++ })
		}()
	}
	wg.Wait()
	return initialized
}
