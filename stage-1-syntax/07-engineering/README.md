# 07. 工程化基础

> 阶段：① 语法精通 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：1.5 天

## 🎯 学习目标

把“写代码”升级为“做工程”——模块、测试、性能、调试样样通。

## 🧩 关键知识点

- module 与 `go.work`、版本语义
- 单元测试：`testing` 包、表驱动测试、子测试
- benchmark 与 `go test -bench`
- lint：`go vet` / `golangci-lint`
- 调试：`dlv` / IDE 断点 / `log` & `slog`
- 性能 profile：`pprof`、CPU / 内存 / 阻塞

## 📦 本章产出

本章提供一组可运行、可测试、可 benchmark 的工程化基础示例：

```text
stage-1-syntax/07-engineering/
├── main.go                  # 组装工程化基础学习报告
├── main_test.go             # 入口报告测试
├── moduleinfo/              # module / go.work / semantic version 摘要
├── calc/                    # 表驱动测试与 benchmark 的纯函数示例
├── quality/                 # go vet / go test / go build / golangci-lint 命令清单
├── debugx/                  # slog 输出与 dlv / IDE 调试入口
├── profile/                 # pprof CPU / memory / blocking 概念与命令提示
└── EXERCISES.md             # 课后练习与验收标准
```

运行示例：

```bash
go run ./stage-1-syntax/07-engineering
```

运行本章测试：

```bash
go test ./stage-1-syntax/07-engineering/...
```

运行 benchmark：

```bash
go test -bench=. ./stage-1-syntax/07-engineering/...
```

与 CI 对齐的本地质量门禁：

```bash
go vet ./...
go test -race -count=1 ./...
go build ./...
golangci-lint run
```

> `golangci-lint` 和 `dlv` 属于外部工具；本章示例会说明命令入口，但单元测试不强制依赖它们已安装。

## 🔗 前置依赖

- 第 06 章

## 📚 推荐扩展阅读

- 《Go 语言高级编程》第 2~3 章
- [Go Testing Documentation](https://pkg.go.dev/testing)
- [golangci-lint](https://golangci-lint.run/)

## ✅ 自测清单

- [ ] 能解释 `go.mod`、module path、Go 版本与依赖声明的关系
- [ ] 能说出 `go.work` 适合多 module 本地协作而不是替代 `go.mod`
- [ ] 能写出一组高质量的表驱动测试和子测试
- [ ] 能用 `go test -bench` 跑出一个真实可比的 benchmark，并理解 `ns/op`
- [ ] 能用 `go vet`、`go test -race`、`go build` 和 `golangci-lint` 做本地质量门禁
- [ ] 能用 `slog` 输出结构化日志辅助定位问题
- [ ] 能用 `dlv test` 或 IDE breakpoint 调试单元测试
- [ ] 能说出 pprof CPU / memory / blocking profile 各自适合定位什么问题
- [ ] 能用 `go test -bench=. -cpuprofile=cpu.out` 与 `go tool pprof` 开始一次性能分析
