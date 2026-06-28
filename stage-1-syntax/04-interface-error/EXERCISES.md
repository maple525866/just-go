# 第 04 章练习：接口、错误与泛型

> 建议先运行 `go test ./stage-1-syntax/04-interface-error/...`，确认示例代码全部通过，再开始练习。

## 练习 1：新增一个隐式实现接口的类型

在 `iface` 包中新增一个 `Video` 类型，并为它实现 `Describe() string` 方法，使它隐式满足 `Describer` 接口。

**验收标准：**

- 为 `Video` 增加至少一个表驱动测试用例，调用 `BuildReport(Video{...})`。
- 测试能够证明 `Video` 不需要显式声明 `implements` 也能作为 `Describer` 使用。
- 运行 `go test ./stage-1-syntax/04-interface-error/iface` 通过。

## 练习 2：扩展 any 分类函数

扩展 `iface.ClassifyAny`，让它能识别 `bool` 类型，并返回形如 `bool:true` 或 `bool:false` 的结果。

**验收标准：**

- 在 `TestClassifyAny` 中新增 `bool` 用例。
- 原有 `string` / `int` / `Describer` / `nil` / `unknown` 用例仍然通过。
- 运行 `go test ./stage-1-syntax/04-interface-error/iface` 通过。

## 练习 3：为错误链增加一层业务上下文

新增一个导出函数，例如 `LoadProfile(name string) error`，它调用 `FindUser` 并再次使用 `fmt.Errorf("load profile: %w", err)` 包装错误。

**验收标准：**

- 对 `LoadProfile("Zoe")` 返回的错误，`errors.Is(err, ErrUserNotFound)` 仍为 `true`。
- 对同一个错误，`errors.As(err, &queryErr)` 仍能提取 `*QueryError`，并读到 `User == "Zoe"`。
- 测试中必须覆盖成功路径和失败路径。

## 练习 4：实现一个泛型 Reduce

在 `generic` 包中实现：

```go
func Reduce[T any, R any](items []T, initial R, fn func(R, T) R) R
```

它从 `initial` 开始遍历切片，把每个元素累积成最终结果。

**验收标准：**

- 用 `Reduce([]int{1,2,3}, 0, ...)` 得到 `6`。
- 用 `Reduce([]string{"Go", "泛型"}, "", ...)` 拼接出一个字符串。
- 使用表驱动 + `t.Run` 编写测试。
- 运行 `go test ./stage-1-syntax/04-interface-error/generic` 通过。

## 练习 5：把示例串到入口报告

选择上面任意一个扩展点，把它加入 `main.go` 的报告输出中。

**验收标准：**

- `go run ./stage-1-syntax/04-interface-error` 的输出包含你新增的示例说明。
- `main_test.go` 中至少有一个断言能覆盖新增输出。
- 运行 `go test ./stage-1-syntax/04-interface-error/...` 通过。
