// Package ptr 演示指针与值/引用语义：
// 值接收者操作的是副本（不改原对象），指针接收者操作的是原对象（原地修改）。
// 这也解释了"何时该用指针接收者"：需要修改接收者，或结构体较大想避免拷贝时。
//
// import 路径：just-go/stage-1-syntax/03-composite-types/ptr
package ptr

// Account 持有一个余额，用于对比值接收者与指针接收者的差异。
type Account struct {
	Balance int
}

// WithBonus 是值接收者方法：在副本上加 bonus 并返回新值，原对象保持不变。
func (a Account) WithBonus(bonus int) Account {
	a.Balance += bonus
	return a
}

// AddBonus 是指针接收者方法：通过指针原地修改原对象的余额。
func (a *Account) AddBonus(bonus int) {
	a.Balance += bonus
}

// Deref 演示取址与解引用：返回指针 p 所指向的值（`*p`）。
func Deref(p *int) int {
	return *p
}

// DoubleInPlace 演示通过指针参数原地修改调用方的变量（取址 `&` + 解引用 `*`）。
func DoubleInPlace(p *int) {
	*p *= 2
}
