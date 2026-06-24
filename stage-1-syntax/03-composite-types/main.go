package main

import (
	"fmt"
	"strings"

	"just-go/stage-1-syntax/03-composite-types/dict"
	"just-go/stage-1-syntax/03-composite-types/model"
	"just-go/stage-1-syntax/03-composite-types/ptr"
	"just-go/stage-1-syntax/03-composite-types/seq"
)

func main() {
	fmt.Println(buildReport("三年二班"))
}

// buildReport 串联 seq / dict / model / ptr 四个子包，组装班级花名册报告正文。
func buildReport(class string) string {
	students := []model.Student{
		{Name: "Ada", Score: 92, Contact: model.Contact{Email: "ada@go.dev"}},
		{Name: "Bob", Score: 78, Contact: model.Contact{Email: "bob@go.dev"}},
		{Name: "Cara", Score: 88, Contact: model.Contact{Email: "cara@go.dev"}},
	}

	scores := make(map[string]int, len(students))
	for _, s := range students {
		scores[s.Name] = s.Score
	}

	lines := []string{
		reportTitle(class),
		"",
		"## 学生（model：struct + 嵌入）",
	}
	for _, s := range students {
		lines = append(lines, "- "+s.Label())
	}

	lines = append(lines,
		"",
		"## 分数表（dict：map + comma-ok）",
		fmt.Sprintf("总分: %d", dict.Total(scores)),
		fmt.Sprintf("及格人数(>=60): %d", dict.CountAtLeast(scores, 60)),
	)
	if score, found := dict.Lookup(scores, "Ada"); found {
		lines = append(lines, fmt.Sprintf("Ada: %d", score))
	}
	if _, found := dict.Lookup(scores, "Zoe"); !found {
		lines = append(lines, "Zoe: (未登记，comma-ok 区分缺失键)")
	}

	lines = append(lines,
		"",
		"## 切片（seq：append 扩容 + 共享底层数组）",
	)
	for i, snap := range seq.GrowSteps(4) {
		lines = append(lines, fmt.Sprintf("append #%d -> len=%d cap=%d", i+1, snap[0], snap[1]))
	}
	base, sub := seq.SubSliceMutationDemo()
	lines = append(lines, fmt.Sprintf("共享底层数组: base=%v sub=%v", base, sub))

	lines = append(lines,
		"",
		"## 指针（ptr：值接收者 vs 指针接收者）",
	)
	acc := ptr.Account{Balance: 100}
	valResult := acc.WithBonus(10)
	balanceAfterValue := acc.Balance
	acc.AddBonus(10)
	lines = append(lines,
		fmt.Sprintf("值接收者返回新值 %d，原对象仍为 %d", valResult.Balance, balanceAfterValue),
		fmt.Sprintf("指针接收者原地修改后原对象为 %d", acc.Balance),
	)

	return strings.Join(lines, "\n")
}

// reportTitle 生成报告标题行（纯函数，便于 main_test 断言）。
func reportTitle(class string) string {
	if class == "" {
		class = "Class"
	}
	return fmt.Sprintf("=== 班级花名册报告：%s ===", class)
}
