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

## 📦 本章产出

本章提供一组可运行、可测试的并发编程示例：

```text
stage-1-syntax/05-concurrency/
├── main.go                  # 组装并发编程学习报告
├── main_test.go             # 入口报告测试
├── goroutine/
│   ├── goroutine.go         # goroutine 启动、等待与结果收集
│   └── goroutine_test.go    # goroutine 表驱动测试
├── channel/
│   ├── channel.go           # 无缓冲/有缓冲 channel、close/range、select timeout
│   └── channel_test.go      # channel 表驱动测试
├── syncx/
│   ├── syncx.go             # Mutex、RWMutex、WaitGroup、Once
│   └── syncx_test.go        # sync 表驱动测试
├── ctx/
│   ├── ctx.go               # context 取消、超时与协作式退出
│   └── ctx_test.go          # context 表驱动测试
├── pitfall/
│   ├── pitfall.go           # data race / goroutine leak / deadlock 风险说明
│   └── pitfall_test.go      # pitfall 表驱动测试
└── EXERCISES.md             # 课后练习与验收标准
```

运行示例：

```bash
go run ./stage-1-syntax/05-concurrency
```

运行本章测试（推荐打开 race detector）：

```bash
go test -race ./stage-1-syntax/05-concurrency/...
```

## 🔗 前置依赖

- 第 04 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 8~9 章
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- 实战工具：`go test -race`

## ✅ 自测清单

- [ ] 能解释 goroutine 启动后为什么需要等待或取消机制
- [ ] 能区分无缓冲 channel 与有缓冲 channel 的通信语义
- [ ] 能正确使用 close + range 读取 channel 中的全部值
- [ ] 能用 `select` 实现 timeout 分支，避免永久阻塞
- [ ] 能用 `sync.WaitGroup` 协调一组并发任务
- [ ] 能用 `Mutex` / `RWMutex` 保护共享状态，并说明何时需要锁
- [ ] 能用 `sync.Once` 保证初始化逻辑只执行一次
- [ ] 能正确使用 `select + context` 实现可取消的并发任务
- [ ] 能识别 data race、goroutine 泄漏、channel deadlock 的成因与规避方式
- [ ] 能用 `go test -race` 检测代码是否存在数据竞争
