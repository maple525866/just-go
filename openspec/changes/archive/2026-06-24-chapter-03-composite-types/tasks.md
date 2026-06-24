## 1. 子包：seq（数组与切片）

- [x] 1.1 创建 `stage-1-syntax/03-composite-types/seq/seq.go`：`package seq`，在顶部注释标注完整 import 路径 `just-go/stage-1-syntax/03-composite-types/seq`，含至少一个未导出辅助标识符演示可见性
- [x] 1.2 在 `seq.go` 实现导出函数演示 `append` 扩容并返回可断言的 `len`/`cap` 变化（如 `GrowDemo(n int) []string` 或返回 `(len, cap)` 序列）
- [x] 1.3 在 `seq.go` 实现导出函数演示切片共享底层数组踩坑（如 `SubSliceMutationDemo`），以返回值而非仅打印体现"改子切片污染原切片"
- [x] 1.4 创建 `stage-1-syntax/03-composite-types/seq/seq_test.go`：表驱动 + `t.Run` 覆盖扩容观察与共享底层数组行为

## 2. 子包：dict（map）

- [x] 2.1 创建 `stage-1-syntax/03-composite-types/dict/dict.go`：`package dict`，标注 import 路径，含注释说明 map 并发写不安全
- [x] 2.2 在 `dict.go` 实现导出函数演示 `comma-ok` 查询（区分"键不存在"与"值为零值"），以及一个基于 map 的统计/汇总纯函数
- [x] 2.3 创建 `stage-1-syntax/03-composite-types/dict/dict_test.go`：表驱动 + `t.Run` 覆盖存在键 / 缺失键 / 零值键场景

## 3. 子包：model（struct）

- [x] 3.1 创建 `stage-1-syntax/03-composite-types/model/model.go`：`package model`，定义含字段标签（如 `json:"..."`）的结构体（如 `Student`）
- [x] 3.2 在 `model.go` 通过嵌入其他结构体（如 `Contact`）演示组合优于继承，并提供导出函数/方法体现字段提升
- [x] 3.3 创建 `stage-1-syntax/03-composite-types/model/model_test.go`：表驱动 + `t.Run` 覆盖结构体构造、嵌入字段访问

## 4. 子包：ptr（指针与值/引用语义）

- [x] 4.1 创建 `stage-1-syntax/03-composite-types/ptr/ptr.go`：`package ptr`，标注 import 路径，顶部注释简述值语义 vs 引用语义
- [x] 4.2 在 `ptr.go` 为同一结构体定义值接收者方法（返回新值、不改原对象）与指针接收者方法（原地修改），并提供取址/解引用演示
- [x] 4.3 创建 `stage-1-syntax/03-composite-types/ptr/ptr_test.go`：表驱动 + `t.Run` 断言值接收者不改原值、指针接收者改原值

## 5. 入口程序

- [x] 5.1 创建 `stage-1-syntax/03-composite-types/main.go`：`package main`，import `seq` / `dict` / `model` / `ptr` 子包，组装并 `fmt.Println` 输出"班级花名册报告"
- [x] 5.2 创建 `stage-1-syntax/03-composite-types/main_test.go`：对入口处可测的纯逻辑（如报告标题/编排）做表驱动测试（不捕获 stdout）

## 6. 学习材料

- [x] 6.1 创建 `stage-1-syntax/03-composite-types/EXERCISES.md`：3~5 道由浅入深的练习题，覆盖切片扩容/共享、map comma-ok、struct 嵌入、值 vs 指针接收者，每题含明确验收标准
- [x] 6.2 更新 `stage-1-syntax/03-composite-types/README.md` 的"📦 本章产出"段落：移除"尚未实现 / 待 OpenSpec change 填充"占位，改列 `.go` 文件清单（main.go / seq/ / dict/ / model/ / ptr/ / *_test.go）+ 运行命令（`go run .`、`go test ./...`）
- [x] 6.3 更新该 README 的"✅ 自测清单"：把检查项改为与本章实际产出对应的可勾选项（对齐 ROADMAP 关键知识点）

## 7. 跨章文档

- [x] 7.1 更新 `docs/glossary.md`：追加本章术语（数组 / 切片 / 容量 cap / map / 结构体 / 字段标签 / 嵌入 / 指针 / 值语义 vs 引用语义 / 指针接收者等），并在"出现章节"列标注 03，保持按英文名字母序

## 8. 验证

- [x] 8.1 在 `stage-1-syntax/03-composite-types/` 执行 `go run .`，确认退出码 0 且打印班级花名册报告
- [x] 8.2 在仓库根目录执行 `go build ./...`，确认退出码 0
- [x] 8.3 在仓库根目录执行 `go test ./stage-1-syntax/03-composite-types/...`，确认全部测试通过
- [x] 8.4 执行 `go vet ./...`，确认无告警
- [x] 8.5 运行 `openspec validate chapter-03-composite-types --strict`，确认本 change 全部产出物合规
