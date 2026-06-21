## Context

`just-go` 第 01、02 章已落地，`go test ./...` 与 `go build ./...` 本地可通过，但 PR 合并前无自动化验证。奠基提案 `bootstrap-learning-repo` 的 design 明确将 CI/CD 列为 Non-Goal，现作为独立基础设施 change 引入。

约束：
- 单 module（`go 1.24`），无第三方 Go 依赖。
- 学习者 Go 水平 L0~L4，lint 规则需学习友好。
- 第 07 章将讲解 golangci-lint，CI 作为前置示范。

## Goals / Non-Goals

**Goals:**
- PR 与 push 到 `main` 时自动运行 test + lint 双 Job。
- 自定义 `.golangci.yml`，启用核心 linter 集，对 `_test.go` 放宽 errcheck。
- 本地命令与 CI 步骤对齐，便于学习者复现。
- 新增 `github-actions-ci` capability spec，定义可验收的 CI 契约。

**Non-Goals:**
- 不配置 branch protection（GitHub 网页手动开启）。
- 不引入 pre-commit hook、Makefile、Docker。
- 不做 Go version matrix。
- 不修改章节教学代码（除非现有代码过不了 lint，届时最小修复）。

## Decisions

### D1：单一 Workflow + 并行 Job

**选择**：`.github/workflows/ci.yml` 含 `test` 与 `lint` 两个并行 Job。

**备选**：拆成两个 workflow / reusable workflow——对当前仓库过度设计。

**理由**：结构清晰、维护成本低，PR 上失败原因一目了然。

### D2：test Job 步骤顺序

`go vet ./...` → `go test -race -count=1 ./...` → `go build ./...`

**理由**：vet 最快失败；race detector 为第 05 章并发内容提前建立习惯；build 兜底编译错误。

### D3：golangci-lint 核心集 + 测试文件 errcheck 豁免

启用：govet、staticcheck、gofmt、gosimple、ineffassign、unused、errcheck、misspell。

不启用：revive、gocyclo、funlen、exhaustive。

**理由**：学习仓库代码量小，风格/复杂度 linter 噪音大于收益。

### D4：Go 版本固定 1.24

与 `go.mod` 一致，不做 matrix。

**理由**：学习路径统一工具链版本，避免「本地过、CI 不过」的版本漂移。

### D5：Action 版本 pin

- `actions/checkout@v4`
- `actions/setup-go@v5`
- `golangci/golangci-lint-action@v7`（pin golangci-lint 版本如 v2.1.6）

**理由**：避免上游 action 或 linter 规则突变导致 CI 无故失败。

## Risks / Trade-offs

| 风险 | 缓解 |
|------|------|
| 现有代码过不了 golangci-lint | 落地前本地跑通；必要时最小修复 |
| Actions 分钟数消耗 | 仅 PR + main push 触发；双 Job 并行约 2~3 min |
| 学习者本地未装 golangci-lint | README 补充安装指引；第 07 章正式讲解 |
| 奠基 Non-Goal 与本次 change 语义冲突 | learning-curriculum spec ADDED Requirement 明确 CI 为后续基础设施演进 |

## Migration Plan

纯增量：新增 `.github/workflows/ci.yml` 与 `.golangci.yml`。合并后首次 push 到 main 即触发 CI。回滚：`git revert` 本 change 提交即可移除 workflow。

## Open Questions

无。
