## 1. golangci-lint 配置

- [ ] 1.1 创建 `.golangci.yml`：启用 govet、staticcheck、gofmt、gosimple、ineffassign、unused、errcheck、misspell；设置 `run.timeout: 5m`、`run.tests: true`；对 `_test.go` 豁免 errcheck
- [ ] 1.2 本地安装 golangci-lint 并执行 `golangci-lint run`，确认现有第 01、02 章代码零 issue；如有违规，做最小修复

## 2. GitHub Actions Workflow

- [ ] 2.1 创建 `.github/workflows/ci.yml`：定义 `pull_request`（目标 `main`）与 `push`（`main`）触发条件
- [ ] 2.2 实现 `test` Job：checkout → setup-go 1.24 → go vet → go test -race -count=1 → go build，超时 10 min
- [ ] 2.3 实现 `lint` Job：checkout → setup-go 1.24 → golangci-lint-action v6（pin 版本），超时 10 min
- [ ] 2.4 确保 `test` 与 `lint` 两个 Job 并行执行（无 depends-on 依赖）

## 3. 文档

- [ ] 3.1 在根目录 `README.md` 追加「本地 CI 对齐」小节：列出四条本地命令（go vet / go test -race / go build / golangci-lint run）

## 4. 验证

- [ ] 4.1 本地执行 `go vet ./...`、`go test -race -count=1 ./...`、`go build ./...`、`golangci-lint run`，全部通过
- [ ] 4.2 push 到 GitHub 后确认 workflow 在 PR 与 main push 上均绿灯
- [ ] 4.3 运行 `openspec validate add-github-actions-ci --strict`，确认本 change 制品合规
