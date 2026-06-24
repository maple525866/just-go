// Package model 演示 struct：字段标签（tag）与通过嵌入实现"组合优于继承"。
//
// import 路径：just-go/stage-1-syntax/03-composite-types/model
package model

import "fmt"

// Contact 是一组可被复用的联系方式字段，作为被嵌入的"零件"。
// 字段标签（如 `json:"email"`）是附在字段上的元信息，常被 encoding/json
// 等库读取以决定序列化时的键名。
type Contact struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Student 通过匿名嵌入 Contact 复用其字段：Go 没有继承，而是用嵌入做组合。
// 嵌入后 Contact 的字段会被"提升"，可像 Student 自身字段一样直接访问
// （如 s.Email 等价于 s.Contact.Email）。
type Student struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Contact        // 匿名嵌入：组合优于继承
}

// Label 利用被提升的嵌入字段 Email，生成一行展示文本，体现字段提升的便利。
func (s Student) Label() string {
	return fmt.Sprintf("%s(%d) <%s>", s.Name, s.Score, s.Email)
}
