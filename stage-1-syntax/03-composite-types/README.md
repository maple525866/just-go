# 03. 复合类型

> 阶段：① 语法精通 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：1 天

## 🎯 学习目标

玩转 array / slice / map / struct / 指针，理解值语义与引用语义。

## 🧩 关键知识点

- 数组与切片：底层数组、`len/cap`、扩容机制、共享底层数组的坑
- map：声明、零值陷阱、并发不安全
- struct：字段标签（tag）、组合（嵌入）
- 指针：取址、解引用、何时用指针接收者
- 值类型 vs 引用类型的传递成本

## 📦 本章产出

**建议阅读顺序**：`seq/` → `dict/` → `model/` → `ptr/` → `main.go`

**示例代码（`.go` 文件清单）：**

| 文件 | 职责 |
| ---- | ---- |
| `main.go` | `package main` 入口，串联四子包输出「班级花名册报告」；含纯函数 `reportTitle` / `buildReport` |
| `main_test.go` | `package main`，表驱动测试标题与报告结构 |
| `seq/seq.go` | 数组 vs 切片、`len`/`cap`、`append` 扩容（`GrowSteps`）、切片共享底层数组踩坑（`SubSliceMutationDemo`）、数组值拷贝（`ArrayValueCopy`） |
| `seq/seq_test.go` | 表驱动 + `t.Run` 测试扩容观察、共享底层数组、数组值语义 |
| `dict/dict.go` | map 的 `comma-ok` 查询（`Lookup`）、汇总（`Total`）、计数（`CountAtLeast`）；注释说明并发不安全 |
| `dict/dict_test.go` | 表驱动测试存在键 / 缺失键 / 零值键、`Total`、`CountAtLeast` |
| `model/model.go` | `Student` 含字段标签（tag），嵌入 `Contact` 演示组合优于继承与字段提升（`Label`） |
| `model/model_test.go` | 表驱动测试嵌入字段提升与 `Label` |
| `ptr/ptr.go` | 值接收者 `WithBonus` vs 指针接收者 `AddBonus`、取址/解引用（`Deref` / `DoubleInPlace`） |
| `ptr/ptr_test.go` | 表驱动测试值/指针接收者语义、解引用 |

子包 import 路径示例：`just-go/stage-1-syntax/03-composite-types/seq`（由 `module just-go` + 目录路径拼成）。

**学习材料：** [`EXERCISES.md`](./EXERCISES.md) —— 5 道由浅入深、含验收标准的练习题。

**运行命令：**

```bash
# 在本章目录运行入口程序
cd stage-1-syntax/03-composite-types
go run .

# 在仓库根目录运行本章测试
go test ./stage-1-syntax/03-composite-types/...
```

## 🔗 前置依赖

- 第 02 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 4 章
- [Go Slices: usage and internals](https://go.dev/blog/slices-intro)

## ✅ 自测清单

- [ ] 能默写 `slice` 扩容的策略，并能在 `seq.GrowSteps` 输出中指出 `cap` 何时增长
- [ ] 能解释为什么 map 不能安全并发写（`fatal error: concurrent map writes`）
- [ ] 能用 `comma-ok` 区分「键不存在」与「值为零值」，避免 map 零值陷阱
- [ ] 能讲清值接收者 vs 指针接收者的选择标准（需修改接收者 / 结构体较大避免拷贝）
- [ ] 能演示一次"切片共享底层数组"踩坑（`seq.SubSliceMutationDemo`）
- [ ] 能说明数组是值类型、切片是引用语义的区别（`seq.ArrayValueCopy`）
- [ ] 能用 struct 嵌入实现"组合优于继承"并解释字段提升（`model.Student` 嵌入 `Contact`）
- [ ] 能在本章目录运行 `go run .` 并看到完整班级花名册报告
- [ ] 能在仓库根运行 `go test ./stage-1-syntax/03-composite-types/...` 并看到全部 PASS
