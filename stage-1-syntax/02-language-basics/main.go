package main

import (
	"fmt"
	"strings"

	"just-go/stage-1-syntax/02-language-basics/control"
	"just-go/stage-1-syntax/02-language-basics/funcs"
	"just-go/stage-1-syntax/02-language-basics/vars"
)

func main() {
	fmt.Println(buildReport("Ada"))
}

// buildReport 串联 vars / control / funcs 三个子包，组装语法基础报告正文。
func buildReport(student string) string {
	scores := []int{92, 78, 88}

	lines := []string{
		reportTitle(student),
		"",
		"## 变量与类型（vars）",
		"零值: " + vars.ZeroValueDemo(),
	}
	for _, line := range vars.DemoSubjects() {
		lines = append(lines, line)
	}

	lines = append(lines,
		"",
		"## 控制流（control）",
		control.Summarize(scores),
	)
	for _, s := range scores {
		lines = append(lines, control.GradeLine(s))
	}

	lines = append(lines,
		"",
		"## 函数（funcs）",
		funcs.MinMaxLine(scores...),
		fmt.Sprintf("average: %.1f", funcs.Average(scores...)),
	)

	grader := funcs.MakeGrader(60)
	lines = append(lines, fmt.Sprintf("grader(78): %s", grader(78)))

	body := strings.Join(lines, "\n")
	return control.RunReport(body)
}

// reportTitle 生成报告标题行（纯函数，便于 main_test 断言）。
func reportTitle(student string) string {
	if student == "" {
		student = "Student"
	}
	return fmt.Sprintf("=== 语法基础报告：%s ===", student)
}
