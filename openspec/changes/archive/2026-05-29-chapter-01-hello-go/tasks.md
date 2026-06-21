## 1. 子包：greeting

- [x] 1.1 创建 `stage-1-syntax/01-hello-go/greeting/greeting.go`：`package greeting`，导出函数 `Greet(name string) string`（name 为空时回退为默认问候对象），内部使用一个未导出的辅助标识符演示可见性
- [x] 1.2 在 `greeting.go` 顶部用注释标注完整 import 路径 `just-go/stage-1-syntax/01-hello-go/greeting` 及其由 `module just-go` + 目录拼成的来历
- [x] 1.3 创建 `stage-1-syntax/01-hello-go/greeting/greeting_test.go`：用 `[]struct{ name, in, want }` 表驱动 + `t.Run` 覆盖正常 / 空字符串两类用例

## 2. 入口程序

- [x] 2.1 创建 `stage-1-syntax/01-hello-go/main.go`：`package main`，import greeting 子包，`main` 中调用 `greeting.Greet(...)` 并 `fmt.Println` 输出
- [x] 2.2 创建 `stage-1-syntax/01-hello-go/main_test.go`：对入口处可测的纯逻辑（如默认名拼装）做一条表驱动测试（不捕获 stdout）

## 3. 学习材料

- [x] 3.1 创建 `stage-1-syntax/01-hello-go/EXERCISES.md`：3~5 道由浅入深的练习题，每题含明确验收标准（如"运行 `go run .` 输出 X"）
- [x] 3.2 更新 `stage-1-syntax/01-hello-go/README.md` 的"📦 本章产出"段落：移除"待 OpenSpec change 填充"占位，改列 `.go` 文件清单（main.go / greeting/ / *_test.go）+ 运行命令（`go run .`、`go test ./...`）
- [x] 3.3 更新该 README 的"✅ 自测清单"：把检查项改为与本章实际产出对应的可勾选项

## 4. 跨章文档

- [x] 4.1 更新 `docs/glossary.md`：追加本章术语（package / module / import path / 可见性 / GOPROXY 等），并在"出现章节"列标注 01

## 5. 验证

- [x] 5.1 在 `stage-1-syntax/01-hello-go/` 执行 `go run .`，确认退出码 0 且打印问候语
- [x] 5.2 在仓库根目录执行 `go build ./...`，确认退出码 0
- [x] 5.3 在仓库根目录执行 `go test ./stage-1-syntax/01-hello-go/...`，确认全部测试通过
- [x] 5.4 执行 `go vet ./...`，确认无告警
- [x] 5.5 运行 `openspec validate chapter-01-hello-go --strict`，确认本 change 4 件产出物合规
