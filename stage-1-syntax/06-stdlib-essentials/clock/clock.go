// Package clock 演示 time：格式化、Duration、Timer/Ticker。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/clock
package clock

import "time"

// Summary 是时间示例返回的结构化结果。
type Summary struct {
	Formatted string
	Minutes   int
}

// FormatSummary 使用 time.Format 和 time.Duration 返回可断言摘要。
func FormatSummary(t time.Time, d time.Duration) Summary {
	return Summary{
		Formatted: t.UTC().Format(time.RFC3339),
		Minutes:   int(d.Minutes()),
	}
}

// CountTicks 使用 time.Ticker 产生有限次数 tick。
func CountTicks(n int, interval time.Duration) int {
	if n <= 0 {
		return 0
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	count := 0
	for count < n {
		<-ticker.C
		count++
	}
	return count
}
