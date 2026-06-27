# Capstone 1 练习：CLI Todo 扩展

> 建议先运行 `go test ./stage-1-syntax/capstone-1-cli-todo/...` 与 `go test -bench=. ./stage-1-syntax/capstone-1-cli-todo/...`，确认基础项目可用。

## 练习 1：增加 update 子命令

新增 `update <id> <title>`，用于修改已有任务标题。

**验收标准：**

- `todo.List` 新增更新标题的方法，并更新 `UpdatedAt`。
- `app.Run` 支持 `update` 子命令。
- 更新不存在 id 时返回可用 `errors.Is` 分类的 `ErrTaskNotFound`。
- 使用表驱动测试覆盖成功、缺少标题、id 不存在三种情况。

## 练习 2：增加只看未完成任务的过滤

为 `list` 增加 `--active` 选项，只显示未完成任务。

**验收标准：**

- 已完成任务不会出现在 `list --active` 输出中。
- 普通 `list` 仍显示全部任务。
- 测试使用临时 JSON 文件和内存 stdout buffer。

## 练习 3：实现导出为纯文本

新增 `export` 子命令，把当前 todo 列表输出为纯文本清单。

**验收标准：**

- 输出包含任务 id、完成状态和标题。
- 空列表输出明确提示。
- 至少为渲染函数增加一个 benchmark。

## 练习 4：改造异步保存为批量 worker

让 `asyncsave.Worker` 支持连续提交多个保存请求，并保证 `Close` 后最后一次提交的数据被持久化。

**验收标准：**

- 连续提交 3 次不同列表，关闭后文件中保留最后一次列表。
- `Close` 后再次 `Submit` 返回 `ErrClosed`。
- `go test -race ./stage-1-syntax/capstone-1-cli-todo/asyncsave` 通过。

## 练习 5：完成阶段答辩笔记

写一段短文，解释这个 CLI Todo 如何综合阶段一 01-07 章知识点。

**验收标准：**

- 覆盖基础语法、struct/slice、error、goroutine/channel、os/json/time、testing/benchmark 六类知识点。
- 说明为什么不引入第三方 CLI 框架。
- 说明如何避免测试污染真实数据文件。
