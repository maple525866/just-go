package main

import (
	"fmt"
	"strings"

	"just-go/stage-1-syntax/04-interface-error/apperr"
	"just-go/stage-1-syntax/04-interface-error/generic"
	"just-go/stage-1-syntax/04-interface-error/iface"
)

func main() {
	fmt.Println(buildReport())
}

// buildReport 串联 iface / apperr / generic 三个子包，组装接口、错误与泛型学习报告。
func buildReport() string {
	book := iface.Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan"}
	lesson := iface.Lesson{Name: "接口与错误", Day: 4}

	_, err := apperr.FindUser("Zoe")
	doubled, longNames, total := generic.Demo()

	lines := []string{
		reportTitle(),
		"",
		"## interface（隐式实现 + 小接口）",
		iface.BuildReport(book).Text,
		iface.BuildReport(lesson).Text,
		"any 分类: " + strings.Join([]string{
			iface.ClassifyAny("hello"),
			iface.ClassifyAny(42),
			iface.ClassifyAny(book),
		}, ", "),
		"",
		"## error（包装 + Is/As）",
		fmt.Sprintf("错误链: %v", err),
		fmt.Sprintf("errors.Is(ErrUserNotFound): %t", apperr.IsUserNotFound(err)),
	}
	if queryErr, ok := apperr.ExtractQueryError(err); ok {
		lines = append(lines, fmt.Sprintf("errors.As(*QueryError): user=%s op=%s", queryErr.User, queryErr.Op))
	}

	lines = append(lines,
		apperr.Summary(),
		"",
		"## generic（类型参数 + 约束）",
		fmt.Sprintf("Map doubled: %v", doubled),
		fmt.Sprintf("Filter long names: %v", longNames),
		fmt.Sprintf("Sum with Number constraint: %d", total),
	)

	return strings.Join(lines, "\n")
}

// reportTitle 生成报告标题行（纯函数，便于 main_test 断言）。
func reportTitle() string {
	return "=== 接口、错误与泛型学习报告 ==="
}
