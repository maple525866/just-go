// Package vars 演示变量、常量、零值与类型转换。
//
// 本包的完整 import 路径为：just-go/stage-1-syntax/02-language-basics/vars
// 它由 go.mod 中的模块前缀（module just-go）拼上本文件所在的目录路径
// （stage-1-syntax/02-language-basics/vars）组成。
package vars

import "fmt"

// 用 iota 定义一组科目编号常量（从 0 递增）。
const (
	subjectMath = iota
	subjectEnglish
	subjectGo
)

// maxScore 是未导出变量，演示包内可见的 var 声明。
var maxScore int = 100

// passScore 用 var 在包级声明及格线。
var passScore int = 60

// subjectLabel 根据 iota 编号返回科目名（未导出辅助函数）。
func subjectLabel(id int) string {
	switch id {
	case subjectMath:
		return "Math"
	case subjectEnglish:
		return "English"
	case subjectGo:
		return "Go"
	default:
		return "Unknown"
	}
}

// FormatScore 格式化单条成绩，含 int → float64 显式转换以计算百分比。
func FormatScore(name string, score int) string {
	// := 在函数内声明局部变量，类型由右值推断。
	pct := float64(score) / float64(maxScore) * 100
	status := "fail"
	if score >= passScore {
		status = "pass"
	}
	return fmt.Sprintf("%s: %d/%.0f (%.1f%%) [%s]", name, score, float64(maxScore), pct, status)
}

// ZeroValueDemo 返回各基本类型零值的字符串表示。
func ZeroValueDemo() string {
	var i int
	var f float64
	var b bool
	var s string
	var r rune
	var by byte
	return fmt.Sprintf(
		"int=%d float=%v bool=%v string=%q rune=%U byte=%d",
		i, f, b, s, r, by,
	)
}

// DemoSubjects 演示 iota 常量组与 := 局部变量，返回三条格式化成绩。
func DemoSubjects() [3]string {
	mathScore := 92    // :=
	englishScore := 78 // :=
	goScore := 88      // :=
	return [3]string{
		FormatScore(subjectLabel(subjectMath), mathScore),
		FormatScore(subjectLabel(subjectEnglish), englishScore),
		FormatScore(subjectLabel(subjectGo), goScore),
	}
}
