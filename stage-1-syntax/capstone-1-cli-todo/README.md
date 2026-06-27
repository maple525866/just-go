# 🚩 Capstone 1: CLI Todo

> 阶段：① 语法精通 · 阶段综合项目 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：3 天

## 🎯 项目目标

综合阶段一所学，做一个**带文件持久化的命令行 Todo**：

- 支持子命令：`add` / `list` / `done` / `delete` / `clear` / `help`
- 数据以 JSON 文件持久化在本地
- 有完整的单元测试与 benchmark
- 与仓库 CI 的 `go test` / `go build` 质量门禁对齐

最终的可执行文件能像 `git` 一样通过子命令调用。

## 🧩 综合应用的章节

- [✓ 01-hello-go] Go 模块与工具链 → 项目骨架
- [✓ 02-language-basics] 函数 / 控制流 → 子命令分发
- [✓ 03-composite-types] struct / slice → 任务数据结构
- [✓ 04-interface-error] 自定义 error → 命令解析错误与任务不存在错误
- [✓ 05-concurrency] goroutine + channel → 受控异步持久化 worker
- [✓ 06-stdlib-essentials] `os` / `encoding/json` / `time` → 文件持久化
- [✓ 07-engineering] 表驱动测试 + benchmark + CI → 质量保障

## 📋 功能清单

```text
stage-1-syntax/capstone-1-cli-todo/
├── main.go              # CLI 入口，委托 app.Run
├── app/                 # 子命令解析、输出渲染、错误映射
├── todo/                # Task/List 领域模型与自定义错误
├── store/               # JSON 文件 Load/Save
├── asyncsave/           # goroutine + channel 异步保存 worker
└── EXERCISES.md         # 扩展练习与验收标准
```

支持命令：

```bash
go run ./stage-1-syntax/capstone-1-cli-todo --help
go run ./stage-1-syntax/capstone-1-cli-todo add "write tests"
go run ./stage-1-syntax/capstone-1-cli-todo list
go run ./stage-1-syntax/capstone-1-cli-todo done 1
go run ./stage-1-syntax/capstone-1-cli-todo delete 1
go run ./stage-1-syntax/capstone-1-cli-todo clear
```

数据文件：

- 默认写入当前目录 `.just-go-todos.json`
- 可通过环境变量指定路径：

```bash
JUST_GO_TODO_FILE=/tmp/todos.json go run ./stage-1-syntax/capstone-1-cli-todo add "learn Go"
```

验证命令：

```bash
go test ./stage-1-syntax/capstone-1-cli-todo/...
go test -bench=. ./stage-1-syntax/capstone-1-cli-todo/...
go run ./stage-1-syntax/capstone-1-cli-todo --help
go test ./...
go build ./...
```

## ✅ 完成标准

- [x] 所有列出的章节知识点至少综合使用过一次
- [x] 代码可运行、有测试、有 README 说明
- [x] 包含 JSON 文件持久化与可配置数据文件路径
- [x] 包含自定义 error、表驱动测试和 benchmark
- [x] 包含受控 goroutine + channel 异步保存示例
- [ ] 阶段答辩（自我口述）：能讲清楚为什么这么设计、踩了哪些坑
