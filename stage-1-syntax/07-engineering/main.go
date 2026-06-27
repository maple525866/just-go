package main

import (
	"fmt"
	"strings"

	"just-go/stage-1-syntax/07-engineering/calc"
	"just-go/stage-1-syntax/07-engineering/debugx"
	"just-go/stage-1-syntax/07-engineering/moduleinfo"
	"just-go/stage-1-syntax/07-engineering/profile"
	"just-go/stage-1-syntax/07-engineering/quality"
)

func main() {
	fmt.Println(buildReport())
}

// buildReport 串联工程化主题子包，组装工程化基础学习报告。
func buildReport() string {
	words := calc.NormalizeWords("Go Test BENCH")
	profiles := profile.Types()
	profileNames := make([]string, 0, len(profiles))
	for _, p := range profiles {
		profileNames = append(profileNames, p.Name)
	}

	return strings.Join([]string{
		reportTitle(),
		"",
		"## module / go.work / semantic version",
		fmt.Sprintf("概念: %s", strings.Join(moduleinfo.Names(), ", ")),
		"",
		"## testing / benchmark",
		fmt.Sprintf("Add(2,3)=%d Fibonacci(10)=%d NormalizeWords=%v", calc.Add(2, 3), calc.Fibonacci(10), words),
		"benchmark: go test -bench=. ./stage-1-syntax/07-engineering/...",
		"",
		"## quality gates",
		strings.Join(quality.Commands(), "\n"),
		"",
		"## debug / slog / dlv",
		strings.TrimSpace(debugx.LogExample("main", 1)),
		strings.Join(debugx.DebugCommands(), "\n"),
		"",
		"## pprof",
		fmt.Sprintf("profiles: %s", strings.Join(profileNames, ", ")),
		strings.Join(profile.Commands(), "\n"),
	}, "\n")
}

// reportTitle 生成报告标题行（纯函数，便于 main_test 断言）。
func reportTitle() string {
	return "=== 工程化基础学习报告 ==="
}
