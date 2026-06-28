// Package moduleinfo 说明 Go module、go.work 与语义化版本。
//
// import 路径：just-go/stage-1-syntax/07-engineering/moduleinfo
package moduleinfo

import "fmt"

// Concept 描述一个工程化概念。
type Concept struct {
	Name    string
	Summary string
}

// Concepts 返回 module、go.work 与 semantic version 的摘要。
func Concepts(moduleName string) []Concept {
	if moduleName == "" {
		moduleName = "just-go"
	}
	return []Concept{
		{Name: "module", Summary: fmt.Sprintf("%s 由 go.mod 声明模块路径、Go 版本与依赖", moduleName)},
		{Name: "go.work", Summary: "go.work 用于把多个 module 组合成一个本地工作区"},
		{Name: "semantic version", Summary: "语义化版本使用 MAJOR.MINOR.PATCH 表达不兼容、功能、修复变更"},
	}
}

// Names 返回概念名称，便于入口报告展示。
func Names() []string {
	concepts := Concepts("just-go")
	names := make([]string, 0, len(concepts))
	for _, concept := range concepts {
		names = append(names, concept.Name)
	}
	return names
}
