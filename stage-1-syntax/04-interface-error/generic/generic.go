// Package generic 演示 Go 泛型：类型参数、约束与可复用切片函数。
//
// import 路径：just-go/stage-1-syntax/04-interface-error/generic
package generic

// Number 是一个类型集约束，限制可参与 Sum 运算的基础数值类型。
type Number interface {
	~int | ~int64 | ~float64
}

// Map 将输入切片中的每个元素转换为另一个类型的元素。
func Map[T any, R any](items []T, fn func(T) R) []R {
	out := make([]R, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	return out
}

// Filter 保留满足条件的元素。
func Filter[T any](items []T, keep func(T) bool) []T {
	out := make([]T, 0, len(items))
	for _, item := range items {
		if keep(item) {
			out = append(out, item)
		}
	}
	return out
}

// Sum 使用 Number 约束演示“只有支持 + 的数值类型才能调用此函数”。
func Sum[T Number](items []T) T {
	var total T
	for _, item := range items {
		total += item
	}
	return total
}

// Demo 返回泛型示例报告中使用的结果。
func Demo() (doubled []int, longNames []string, total int) {
	doubled = Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
	longNames = Filter([]string{"Go", "error", "generic"}, func(s string) bool { return len(s) > 2 })
	total = Sum([]int{10, 20, 30})
	return doubled, longNames, total
}
