# github-actions-ci Specification

## Purpose
定义 just-go 仓库级 GitHub Actions CI 契约——PR 与 main push 触发、test/lint 双 Job 并行、golangci-lint 自定义配置与本地命令对齐。

## Requirements

### Requirement: 仓库 SHALL 提供 GitHub Actions CI workflow

仓库 SHALL 在 `.github/workflows/` 下提供 CI workflow，在 `pull_request`（目标分支 `main`）与 `push`（分支 `main`）时自动触发，且 MUST 包含 `test` 与 `lint` 两个并行 Job。

#### Scenario: PR 触发 CI

- **WHEN** 有人向 `main` 分支提交 pull request
- **THEN** GitHub Actions MUST 自动运行 CI workflow，且 `test` 与 `lint` 两个 Job MUST 并行执行

#### Scenario: main push 触发 CI

- **WHEN** 有人 push 提交到 `main` 分支
- **THEN** GitHub Actions MUST 自动运行 CI workflow

### Requirement: test Job SHALL 执行 vet、race test 与 build

CI 的 `test` Job MUST 依次执行 `go vet ./...`、`go test -race -count=1 ./...`、`go build ./...`，且 MUST 使用与 `go.mod` 一致的 Go 版本（当前为 1.24）。

#### Scenario: test Job 全部通过

- **WHEN** 仓库中所有包通过 vet、race test 与 build
- **THEN** `test` Job MUST 以成功状态结束

#### Scenario: test Job 检测到失败

- **WHEN** 任一 `go vet`、`go test -race` 或 `go build` 命令失败
- **THEN** `test` Job MUST 以失败状态结束，且整个 workflow MUST 标记为失败

### Requirement: lint Job SHALL 使用 golangci-lint 与自定义配置

CI 的 `lint` Job MUST 使用 `golangci-lint run` 对仓库执行静态检查，且 MUST 读取仓库根目录的 `.golangci.yml` 配置文件。

#### Scenario: lint Job 零 issue 通过

- **WHEN** golangci-lint 对仓库报告零 issue
- **THEN** `lint` Job MUST 以成功状态结束

#### Scenario: lint Job 检测到违规

- **WHEN** golangci-lint 报告至少一个 issue
- **THEN** `lint` Job MUST 以失败状态结束，且整个 workflow MUST 标记为失败

### Requirement: 仓库 SHALL 提供 .golangci.yml 配置文件

仓库根目录 SHALL 包含 `.golangci.yml`，MUST 启用至少以下 linter：`govet`、`staticcheck`、`gofmt`、`gosimple`、`ineffassign`、`unused`、`errcheck`、`misspell`，且 MUST 对 `_test.go` 文件豁免 `errcheck` 规则。

#### Scenario: 配置文件存在且可被读取

- **WHEN** 开发者在仓库根目录执行 `golangci-lint run`
- **THEN** 命令 MUST 成功读取 `.golangci.yml` 并执行检查

#### Scenario: 测试文件 errcheck 豁免生效

- **WHEN** `_test.go` 文件中存在被 errcheck 标记的写法（如 `_ = fn()`）
- **THEN** golangci-lint MUST NOT 因此对该测试文件报告 errcheck issue
