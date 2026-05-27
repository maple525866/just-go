# 🚩 Capstone 1: CLI Todo

> 阶段：① 语法精通 · 阶段综合项目 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：3 天

## 🎯 项目目标

综合阶段一所学，做一个**带文件持久化的命令行 Todo**：

- 支持子命令：`add` / `list` / `done` / `delete` / `clear`
- 数据以 JSON 文件持久化在本地
- 有完整的单元测试
- 配置好 `golangci-lint` 并通过

最终的可执行文件应该能像 `git` 一样优雅地被命令行调用。

## 🧩 综合应用的章节

- [✓ 01-hello-go] Go 模块与工具链 → 项目骨架
- [✓ 02-language-basics] 函数 / 控制流 → 子命令分发
- [✓ 03-composite-types] struct / slice → 任务数据结构
- [✓ 04-interface-error] 自定义 error → 命令解析错误
- [✓ 05-concurrency] goroutine + channel → 异步持久化（可选）
- [✓ 06-stdlib-essentials] `os` / `encoding/json` / `time` → 文件持久化
- [✓ 07-engineering] 表驱动测试 + benchmark + lint → 质量保障

## 📋 功能清单（待 OpenSpec change 填充）

> ⚠️ 当前项目尚未实现。
>
> 请通过 `/opsx-propose capstone-1-cli-todo` 来启动本项目。

## ✅ 完成标准（落地后填充）

- [ ] 所有列出的章节知识点至少综合使用过一次
- [ ] 代码可运行、有测试（覆盖率 ≥ 70%）、有 README 说明
- [ ] 通过 `golangci-lint run` 无 warning
- [ ] 阶段答辩（自我口述）：能讲清楚为什么这么设计、踩了哪些坑
