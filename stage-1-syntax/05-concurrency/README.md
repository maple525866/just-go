# 05. 并发编程

> 阶段：① 语法精通 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

## 🎯 学习目标

用 goroutine + channel + sync + context 写出**正确的**并发程序。

## 🧩 关键知识点

- goroutine 启动与生命周期、调度概览（GMP 浅层）
- channel：有缓冲 / 无缓冲、关闭语义、`select`、超时模式
- `sync` 包：`Mutex` / `RWMutex` / `WaitGroup` / `Once`
- `context.Context`：取消、超时、值传递
- 常见并发坑：data race、goroutine 泄漏、channel 死锁

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-05-concurrency` 来落地本章内容。

## 🔗 前置依赖

- 第 04 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 8~9 章
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- 实战工具：`go test -race`

## ✅ 自测清单（落地后填充）

- [ ] 能正确使用 `select + context` 实现可取消的并发任务
- [ ] 能识别并修复一个 goroutine 泄漏
- [ ] 能用 `sync.WaitGroup` 协调一组并发任务
- [ ] 能解释为什么"不要通过共享内存来通信，要通过通信来共享内存"
- [ ] 能用 `go test -race` 检测一段代码的数据竞争
