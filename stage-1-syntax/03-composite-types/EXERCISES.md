# 03. 复合类型 — 练习题

> 以下练习只用本章已学概念：数组 / 切片（`len`/`cap`/`append`/共享底层数组）、map（`comma-ok`/零值）、struct（tag/嵌入）、指针（值 vs 指针接收者）。
> 在 `stage-1-syntax/03-composite-types/` 目录下完成并验证。

## 练习 1：观察切片扩容（⭐）

阅读 `seq/seq.go` 的 `GrowSteps`，在本目录运行 `go run .`，记录「切片（seq）」段落里每次 `append` 后的 `len` 与 `cap`。

**验收标准**：能说出一处 `cap` 翻倍（增长）发生在第几次 `append`；用一句话解释为什么 `cap` 不是每次都变。

## 练习 2：复现「共享底层数组」踩坑（⭐⭐）

在 `seq/seq_test.go` 新增一个子测试：构造 `base := []int{10, 20, 30}`，切出 `sub := base[:2]`，对 `sub[0]` 赋新值后断言 `base[0]` 也随之改变。

**验收标准**：运行 `go test ./stage-1-syntax/03-composite-types/seq/...` 全部 PASS，且测试能证明 `base` 被 `sub` 的修改污染。

## 练习 3：用 comma-ok 避免零值陷阱（⭐⭐）

在 `dict/dict.go` 新增导出函数 `LookupOrDefault(scores map[string]int, name string, def int) int`：键存在时返回其分数，缺失时返回 `def`（必须用 `comma-ok` 判断，而不是直接 `scores[name]`）。

**验收标准**：新增表驱动测试覆盖「存在键」「缺失键返回 def」「值为 0 的键返回 0 而非 def」三种用例，`go test ./stage-1-syntax/03-composite-types/dict/...` 全部 PASS。

## 练习 4：用嵌入实现组合（⭐⭐）

在 `model/model.go` 新增一个结构体 `Address{ City, Street string }`，把它嵌入 `Student`，并新增方法 `Student.FullAddress() string` 返回 `City + Street`（演示嵌入字段提升）。

**验收标准**：新增测试构造带 `Address` 的 `Student`，断言可直接通过 `s.City` 访问被提升的字段，且 `FullAddress()` 返回拼接结果；`go test ./stage-1-syntax/03-composite-types/model/...` 全部 PASS。

## 练习 5：辨析值接收者 vs 指针接收者（⭐⭐⭐）

在 `ptr/ptr.go` 新增**值接收者**方法 `Account.TryWithdrawValue(n int) Account`（在副本上扣减并返回新值）与**指针接收者**方法 `(*Account).Withdraw(n int)`（原地扣减）。

**验收标准**：新增测试断言「调用值接收者方法后原 `Account` 余额不变」「调用指针接收者方法后原 `Account` 余额减少」；并用一句注释写下你对「何时该用指针接收者」的理解（需要修改接收者 / 结构体较大避免拷贝）。`go test ./stage-1-syntax/03-composite-types/ptr/...` 全部 PASS。
