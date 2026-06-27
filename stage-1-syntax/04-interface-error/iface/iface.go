// Package iface 演示 Go 接口：隐式实现、小接口、any、类型断言与 type switch。
//
// import 路径：just-go/stage-1-syntax/04-interface-error/iface
package iface

import "fmt"

// Describer 是一个小接口：只描述调用方真正需要的方法。
type Describer interface {
	Describe() string
}

// Report 是接口示例返回的具体结构体。
type Report struct {
	Kind string
	Text string
}

// Book 没有显式声明实现 Describer；只要方法集匹配，就隐式满足接口。
type Book struct {
	Title  string
	Author string
}

// Describe 返回书籍说明。
func (b Book) Describe() string {
	if b.Author == "" {
		return b.Title
	}
	return fmt.Sprintf("%s by %s", b.Title, b.Author)
}

// Lesson 也是 Describer 的隐式实现，用于展示同一个小接口可接收不同具体类型。
type Lesson struct {
	Name string
	Day  int
}

// Describe 返回课程说明。
func (l Lesson) Describe() string {
	return fmt.Sprintf("Day %d: %s", l.Day, l.Name)
}

// BuildReport 接受小接口，返回具体结构体；这体现“接受接口，返回结构体”的常见设计原则。
func BuildReport(d Describer) Report {
	return Report{Kind: "describer", Text: d.Describe()}
}

// ClassifyAny 使用 type switch 分析 any 值；any 是 interface{} 的别名。
func ClassifyAny(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case string:
		return "string:" + x
	case int:
		return fmt.Sprintf("int:%d", x)
	case Describer:
		return "describer:" + x.Describe()
	default:
		return fmt.Sprintf("unknown:%T", v)
	}
}

// IsDescriber 使用类型断言判断 any 值是否满足 Describer。
func IsDescriber(v any) bool {
	_, ok := v.(Describer)
	return ok
}
