## ADDED Requirements

### Requirement: 仓库 SHALL 提供 GitHub Actions CI 质量门禁

仓库 SHALL 通过 GitHub Actions 在 PR 与 main push 时自动验证代码质量（test + lint），作为合并前的客观门禁。此要求是对奠基阶段「不配置 CI/CD」非目标的后续演进，由独立 OpenSpec change（如 `add-github-actions-ci`）引入。

#### Scenario: CI workflow 文件存在

- **WHEN** 阅读者列出 `.github/workflows/` 目录
- **THEN** 该目录 MUST 至少包含一个 CI workflow 文件（如 `ci.yml`）

#### Scenario: CI 与本地命令对齐

- **WHEN** 开发者在仓库根目录依次执行 `go vet ./...`、`go test -race -count=1 ./...`、`go build ./...`、`golangci-lint run`
- **THEN** 上述命令的结果 MUST 与 CI workflow 中对应步骤的通过/失败语义一致
