// Package seq 演示数组与切片：len/cap、append 扩容，以及切片共享底层数组的踩坑。
//
// import 路径：just-go/stage-1-syntax/03-composite-types/seq
package seq

// snapshot 返回切片当前的 [len, cap]（未导出辅助函数，演示包内可见性）。
func snapshot(s []int) [2]int {
	return [2]int{len(s), cap(s)}
}

// GrowSteps 从空切片开始连续 append n 次，返回每次 append 之后的 [len, cap] 快照，
// 用于直观观察切片随元素增加而触发的扩容（cap 通常按倍数增长）。
func GrowSteps(n int) [][2]int {
	steps := make([][2]int, 0, n)
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, i)
		steps = append(steps, snapshot(s))
	}
	return steps
}

// SubSliceMutationDemo 演示切片共享底层数组的踩坑：
// sub 是从 base 切出的子切片，二者共享同一底层数组，
// 因此修改 sub[0] 会同时改动 base[1]。返回修改后的 base 与 sub 以便断言"共享"这一事实。
func SubSliceMutationDemo() (base, sub []int) {
	base = []int{1, 2, 3, 4, 5}
	sub = base[1:3] // 共享 base 的底层数组
	sub[0] = 99     // 等价于修改 base[1]
	return base, sub
}

// ArrayValueCopy 演示数组是值类型：把数组赋值给另一个变量是整份拷贝，
// 修改副本不会影响原数组（与切片的引用语义形成对比）。
func ArrayValueCopy() (original, modified [3]int) {
	original = [3]int{1, 2, 3}
	modified = original // 整份拷贝
	modified[0] = 99
	return original, modified
}
