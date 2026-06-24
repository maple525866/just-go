// Package dict 演示 map：声明与读写、comma-ok 查询缺失键、零值陷阱。
//
// import 路径：just-go/stage-1-syntax/03-composite-types/dict
//
// 注意：map 并发写不安全！多个 goroutine 同时写同一个 map 会触发运行时
// `fatal error: concurrent map writes`。需并发访问时应配合 sync.Mutex 保护，
// 或改用 sync.Map（详见第 05 章并发）。本章不引入并发，仅在此提示该陷阱。
package dict

// Lookup 以 comma-ok 形式查询某人的分数：
// found 用来区分"键不存在"与"键存在但值恰为零值 0"，
// 从而避免把缺失键误判成 0 分（这正是 map 的零值陷阱）。
func Lookup(scores map[string]int, name string) (score int, found bool) {
	score, found = scores[name]
	return
}

// Total 遍历 map 汇总所有分数（基于 map 的纯函数，便于单元测试断言）。
func Total(scores map[string]int) int {
	sum := 0
	for _, v := range scores {
		sum += v
	}
	return sum
}

// CountAtLeast 统计分数不低于 threshold 的人数。
func CountAtLeast(scores map[string]int, threshold int) int {
	count := 0
	for _, v := range scores {
		if v >= threshold {
			count++
		}
	}
	return count
}
