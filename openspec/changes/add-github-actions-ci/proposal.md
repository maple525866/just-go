## Why

第 01、02 章已落地可运行代码与测试，但仓库仍无自动化质量门禁——PR 合并前无法自动验证 `go test`、`go build` 与静态检查是否通过。引入 GitHub Actions CI 可在合并前捕获回归，并为第 07 章「工程化基础」中的 lint 主题提供前置示范。

## What Changes

- **新增** `.github/workflows/ci.yml`：双 Job 并行 workflow（test + lint），在 PR 与 push 到 `main` 时触发。
- **新增** `.golangci.yml`：自定义 golangci-lint 规则集，平衡学习友好与工程质量。
- **新增** capability `github-actions-ci`：定义 CI 行为契约（触发条件、检查项、失败语义）。
- **修改** `learning-curriculum`：追加「仓库 SHALL 提供 GitHub Actions CI 门禁」要求，与奠基提案中「不配置 CI/CD」的非目标解耦（基础设施演进）。

## Capabilities

### New Capabilities

- `github-actions-ci`：定义仓库级 GitHub Actions CI 契约——PR 与 main push 触发、test job（go vet / go test -race / go build）与 lint job（golangci-lint）并行执行、任一失败则 workflow 失败。

### Modified Capabilities

- `learning-curriculum`：新增 Requirement，要求仓库提供 GitHub Actions CI 作为代码质量门禁；更新奠基阶段「不配置 CI/CD」的表述，明确 CI 作为后续基础设施 change 引入。

## Impact

- **基础设施**：新增 `.github/` 与 `.golangci.yml`，不影响现有 Go 源码逻辑。
- **依赖**：不引入 Go 第三方依赖；CI 使用 GitHub Actions 官方/社区 action（`actions/checkout`、`actions/setup-go`、`golangci/golangci-lint-action`）。
- **学习者体验**：PR 提交后自动获得红/绿反馈；README 可补充本地对齐命令。
- **Actions 消耗**：每次 PR 与 main push 触发 2 个并行 Job，预计 2~3 分钟/次。
