// Package greeting 演示「包拆分 + import 路径 + 可见性」三件事。
//
// 本包的完整 import 路径为：just-go/stage-1-syntax/01-hello-go/greeting
// 它由 go.mod 中的模块前缀（module just-go）拼上本文件所在的目录路径
// （stage-1-syntax/01-hello-go/greeting）组成 —— 这正是 Go 解析 import 的规则。
package greeting

import "fmt"

// defaultName 是未导出标识符（首字母小写），只在本包内可见，
// 跨包无法访问；以此演示「大小写决定可见性」这一 Go 核心约定。
const defaultName = "Gopher"

// Greet 是导出函数（首字母大写），可被其他包调用。
// name 为空时回退到 defaultName，保证总能返回一句完整问候语。
func Greet(name string) string {
	if name == "" {
		name = defaultName
	}
	return fmt.Sprintf("Hello, %s!", name)
}
