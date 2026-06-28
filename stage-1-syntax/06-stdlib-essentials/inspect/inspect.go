// Package inspect 演示 reflect 的只读元数据检查。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/inspect
package inspect

import "reflect"

// Field 描述结构体字段的只读元信息。
type Field struct {
	Name string
	Type string
	JSON string
}

// DescribeStruct 使用 reflect 读取类型名、字段名、字段类型和 json tag。
func DescribeStruct(v any) (string, []Field) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return t.String(), nil
	}

	fields := make([]Field, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, Field{
			Name: field.Name,
			Type: field.Type.String(),
			JSON: field.Tag.Get("json"),
		})
	}
	return t.Name(), fields
}
