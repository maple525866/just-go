# GitHub Actions CI 设计文档

> 日期：2026-06-21  
> Change：`add-github-actions-ci`  
> 状态：已批准

## 背景

`just-go` 是 Go 学习仓库，采用 OpenSpec 章节式工作流。奠基提案刻意排除了 CI/CD，但第 01、02 章已落地可运行代码与测试。为在 PR 合并前自动验证代码质量，需引入 GitHub Actions CI。

## 目标

- 在 PR 与 push 到 `main` 时自动运行完整工程化检查。
- 提供自定义 `.golangci.yml`，平衡「抓真问题」与「学习友好」。
- 作为第 07 章「工程化基础」的前置示范。

## 非目标

- 不配置 branch protection（需在 GitHub 网页手动开启）。
- 不引入 pre-commit hook、Makefile、Docker。
- 不做 Go version matrix。
- 不修改章节教学代码（除非现有代码过不了 lint，届时最小修复）。

## 架构

### 方案选择

采用**单一 Workflow + 并行 Job**（方案 A）：

```text
.github/workflows/ci.yml     ← 唯一 workflow
.golangci.yml                ← 自定义 lint 规则
```

### 触发条件

- `pull_request` → 目标分支 `main`
- `push` → 分支 `main`

### Go 版本

固定 `1.24`（与 `go.mod` 一致）。

### Job 结构

| Job | 步骤 | 超时 |
|-----|------|------|
| `test` | checkout → setup-go → go vet → go test -race -count=1 → go build | 10 min |
| `lint` | checkout → setup-go → golangci-lint-action | 10 min |

两个 Job 并行，任一失败则 workflow 失败。

## `.golangci.yml` 规则策略

### 启用的 linter

| Linter | 作用 |
|--------|------|
| `govet` | 官方静态分析 |
| `staticcheck` | 深度 bug / 性能 / 风格问题 |
| `gofmt` | 格式一致性 |
| `gosimple` | 可简化写法 |
| `ineffassign` | 无效赋值 |
| `unused` | 未使用代码 |
| `errcheck` | 忽略 error 返回值 |
| `misspell` | 英文拼写 |

### 不启用的 linter

`revive`、`gocyclo`、`funlen`、`exhaustive` 等——学习阶段代码量小，过早引入增加噪音。

### 学习友好配置

```yaml
run:
  timeout: 5m
  tests: true

issues:
  max-issues-per-linter: 0
  exclude-rules:
    - path: _test\.go
      linters: [errcheck]
```

### 版本锁定

`golangci-lint-action` 使用 v6，pin 到具体 minor 版本。

## 本地对齐命令

```bash
go vet ./...
go test -race -count=1 ./...
go build ./...
golangci-lint run
```

## OpenSpec 变更范围

**Change 名称**：`add-github-actions-ci`

**产出文件**：

| 文件 | 内容 |
|------|------|
| `.github/workflows/ci.yml` | 双 Job workflow |
| `.golangci.yml` | 自定义 lint 规则 |
| `openspec/changes/add-github-actions-ci/` | proposal / design / tasks / spec delta |

## 验收标准

1. 本地 `golangci-lint run` 对现有第 01、02 章代码零 issue。
2. push 到 GitHub 后，workflow 在 PR 和 main push 上均绿灯。
3. 故意引入 lint 违规 → CI lint job 失败。
4. `go test -race ./...` 在 CI 环境通过。

## 与第 07 章的关系

CI 先引入 golangci-lint 作为工程化前置示范；第 07 章落地时可在此基础上讲解规则含义，并教学习者本地运行。
