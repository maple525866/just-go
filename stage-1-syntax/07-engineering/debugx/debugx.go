// Package debugx 演示 slog 日志输出，并提供 dlv / IDE 调试入口摘要。
//
// import 路径：just-go/stage-1-syntax/07-engineering/debugx
package debugx

import (
	"bytes"
	"log/slog"
)

// LogExample 使用 slog 写入可断言的文本日志。
func LogExample(component string, step int) string {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	}}))
	logger.Info("debug checkpoint", "component", component, "step", step)
	return buf.String()
}

// DebugCommands 返回常见调试入口。
func DebugCommands() []string {
	return []string{
		"dlv test ./stage-1-syntax/07-engineering/calc",
		"IDE breakpoint: set breakpoint in calc.Fibonacci and run package tests",
		"slog: add structured key/value fields near suspicious branches",
	}
}
