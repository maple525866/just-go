package main

import (
	"fmt"
	"strings"
	"time"

	"just-go/stage-1-syntax/05-concurrency/channel"
	"just-go/stage-1-syntax/05-concurrency/ctx"
	"just-go/stage-1-syntax/05-concurrency/goroutine"
	"just-go/stage-1-syntax/05-concurrency/pitfall"
	"just-go/stage-1-syntax/05-concurrency/syncx"
)

func main() {
	fmt.Println(buildReport())
}

// buildReport 串联各并发主题子包，组装并发编程学习报告。
func buildReport() string {
	squares := goroutine.RunSquares([]int{2, 3, 4})
	board := syncx.NewScoreBoard()
	board.Set("Ada", 95)
	score, found := board.Get("Ada")

	lines := []string{
		reportTitle(),
		"",
		"## goroutine（启动 + 等待 + 收集）",
		fmt.Sprintf("生命周期: %s", strings.Join(goroutine.LifecycleSummary(), " -> ")),
		fmt.Sprintf("平方结果: %v", squares),
		"",
		"## channel（通信 + close + select）",
		fmt.Sprintf("无缓冲通信: %s", channel.PingPong("ping")),
		fmt.Sprintf("有缓冲队列: %v", channel.BufferedQueue([]int{1, 2, 3})),
		fmt.Sprintf("close + range: %v", channel.CloseAndRange([]string{"a", "b"})),
		fmt.Sprintf("select timeout: %s", channel.ReceiveWithTimeout(5*time.Millisecond, time.Millisecond)),
		"",
		"## sync（Mutex + RWMutex + Once）",
		fmt.Sprintf("Mutex 计数: %d", syncx.CountWithMutex(5)),
		fmt.Sprintf("RWMutex 读取: Ada=%d found=%t", score, found),
		fmt.Sprintf("Once 初始化次数: %d", syncx.InitOnceConcurrently(5)),
		"",
		"## context（取消 + 超时）",
		fmt.Sprintf("取消结果: %s", ctx.CancellationDemo()),
		fmt.Sprintf("超时结果: %t", ctx.TimeoutDemo(time.Millisecond)),
		"",
		"## 并发坑（安全说明）",
		fmt.Sprintf("常见风险: %s", strings.Join(pitfall.Names(), ", ")),
	}
	return strings.Join(lines, "\n")
}

// reportTitle 生成报告标题行（纯函数，便于 main_test 断言）。
func reportTitle() string {
	return "=== 并发编程学习报告 ==="
}
