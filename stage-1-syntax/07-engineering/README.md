# 07. 工程化基础

> 阶段：① 语法精通 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：1.5 天

## 🎯 学习目标

把"写代码"升级为"做工程"——模块、测试、性能、调试样样通。

## 🧩 关键知识点

- module 与 `go.work`、版本语义
- 单元测试：`testing` 包、表驱动测试、子测试
- benchmark 与 `go test -bench`
- lint：`go vet` / `golangci-lint`
- 调试：`dlv` / IDE 断点 / `log` & `slog`
- 性能 profile：`pprof`、CPU / 内存 / 阻塞

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-07-engineering` 来落地本章内容。

## 🔗 前置依赖

- 第 06 章

## 📚 推荐扩展阅读

- 《Go 语言高级编程》第 2~3 章
- [Go Testing Documentation](https://pkg.go.dev/testing)
- [golangci-lint](https://golangci-lint.run/)

## ✅ 自测清单（落地后填充）

- [ ] 能写出一组高质量的表驱动测试
- [ ] 能用 `go test -bench` 跑出一个真实可比的 benchmark
- [ ] 能用 `pprof` 找到一段 CPU 热点并优化
- [ ] 能配置 `golangci-lint` 并在本地一键跑通
- [ ] 能用 `dlv` 进行远程 / 本地调试
