# 01. Hello, Go

> 阶段：① 语法精通 ｜ 难度：⭐☆☆☆☆ ｜ 预计耗时：3 小时

## 🎯 学习目标

装好 Go 工具链，跑通第一个程序，理解 `go run / build / mod` 三件套与项目结构。

## 🧩 关键知识点

- Go 安装与环境变量（GOPATH / GOROOT / GOPROXY）
- `go mod init / tidy / run / build / test` 命令
- 包（package）、import 路径、`main` 函数
- VS Code / GoLand 基础配置与调试器

## 📦 本章产出

**示例代码（`.go` 文件清单）：**

| 文件 | 职责 |
| ---- | ---- |
| `main.go` | `package main` 入口，调用 `greeting.Greet` 并 `fmt.Println` 输出；含纯函数 `resolveName` 决定问候对象 |
| `main_test.go` | `package main`，表驱动测试 `resolveName`（不捕获 stdout） |
| `greeting/greeting.go` | `package greeting`，导出 `Greet(name)`，含未导出常量 `defaultName` 演示可见性 |
| `greeting/greeting_test.go` | `package greeting`，表驱动 + `t.Run` 测试 `Greet` |

`greeting` 子包的完整 import 路径为 `just-go/stage-1-syntax/01-hello-go/greeting`，由 `go.mod` 的 `module just-go` 前缀拼接目录路径而成。

**学习材料：** [`EXERCISES.md`](./EXERCISES.md) —— 5 道由浅入深、含验收标准的练习题。

**运行命令：**

```bash
# 在本章目录运行入口程序
cd stage-1-syntax/01-hello-go
go run .

# 在仓库根目录运行本章测试
go test ./stage-1-syntax/01-hello-go/...
```

## 🔗 前置依赖

无（这是整条学习路径的起点）。

## 📚 推荐扩展阅读

- [Tour of Go § Welcome](https://go.dev/tour/welcome/1)
- [Effective Go § Introduction](https://go.dev/doc/effective_go)
- [How to Write Go Code](https://go.dev/doc/code)

## ✅ 自测清单

- [ ] 能在 `stage-1-syntax/01-hello-go/` 运行 `go run .` 并看到问候语输出
- [ ] 能解释 `main` 包与 `greeting` 包的 import 路径（`just-go` 前缀 + 目录路径）如何拼成
- [ ] 能说出导出标识符（首字母大写）与未导出标识符（首字母小写）的区别
- [ ] 能在仓库根运行 `go test ./stage-1-syntax/01-hello-go/...` 并看到全部 PASS
