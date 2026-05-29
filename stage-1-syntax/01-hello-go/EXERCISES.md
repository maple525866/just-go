# 01. Hello, Go — 练习题

> 以下练习只用本章已学概念：`package` / `import` / `main`、`go run/build/test`、导出与未导出标识符（大小写可见性）。
> 在 `stage-1-syntax/01-hello-go/` 目录下完成并验证。

## 练习 1：改一句问候（⭐）

修改 `main.go`，让程序问候你自己的名字（例如 `Ada`）。

**验收标准**：在本目录运行 `go run .`，输出为 `Hello, Ada!`。

## 练习 2：理解空值回退（⭐）

不改 `greeting` 包，只在 `main.go` 里把传给 `greeting.Greet` 的名字改成空字符串 `""`。

**验收标准**：运行 `go run .`，输出为 `Hello, Gopher!`（说明空名时回退到了 `greeting` 包内的默认问候对象）。

## 练习 3：新增一个导出函数（⭐⭐）

在 `greeting/greeting.go` 中新增一个导出函数 `Shout(name string) string`，返回全大写的问候语（提示：用 `strings.ToUpper` 包裹 `Greet` 的结果，并 `import "strings"`）。在 `main.go` 中调用它并打印。

**验收标准**：运行 `go run .`，能看到一行形如 `HELLO, ADA!` 的全大写输出；运行 `go build ./...` 退出码为 0。

## 练习 4：给新函数补一条表驱动测试（⭐⭐）

在 `greeting/greeting_test.go` 中为练习 3 的 `Shout` 增加一个表驱动子测试，至少覆盖「普通名字」与「空字符串」两种用例。

**验收标准**：在仓库根目录运行 `go test ./stage-1-syntax/01-hello-go/...`，全部测试 `PASS`。

## 练习 5：验证未导出标识符不可跨包访问（⭐⭐⭐）

在 `main.go` 中尝试直接引用 `greeting.defaultName`（未导出常量），先观察编译错误，再删除这行代码恢复可编译状态。

**验收标准**：加入该行后运行 `go build ./...` 会报错（提示 `cannot refer to unexported name`），删除后 `go build ./...` 再次退出码 0。请用一句话写下你对「首字母大小写决定可见性」的理解。
