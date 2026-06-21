// Package control 演示 if / for / switch / defer 控制流。
package control

import (
	"fmt"
	"strings"
)

const reportFooter = "--- end of report ---"

// LetterGrade 用 switch 把分数映射为 A~F 等级。
func LetterGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

// Summarize 汇总分数列表：if 初始化语句处理空输入；for 经典三段式求和；for range 拼接明细。
func Summarize(scores []int) string {
	if n := len(scores); n == 0 {
		return "count=0 sum=0 avg=0"
	}

	sum := 0
	for i := 0; i < len(scores); i++ {
		sum += scores[i]
	}

	parts := make([]string, 0, len(scores))
	for _, s := range scores {
		parts = append(parts, fmt.Sprintf("%d", s))
	}

	avg := float64(sum) / float64(len(scores))
	return fmt.Sprintf("count=%d sum=%d avg=%.1f detail=[%s]",
		len(scores), sum, avg, strings.Join(parts, ","))
}

// RunReport 在报告正文末尾 defer 追加页脚；使用命名返回值以便 defer 修改最终返回字符串。
func RunReport(body string) (report string) {
	report = body
	defer func() {
		report += "\n" + reportFooter
	}()
	return report
}

// GradeLine 把单条分数格式化为「分数 → 等级」一行。
func GradeLine(score int) string {
	return fmt.Sprintf("%d → %s", score, LetterGrade(score))
}
