## Context

`stage-1-syntax/04-interface-error/` 当前只有占位 README 与 `.gitkeep`。前 3 章已形成相对稳定的章节落地风格：章节根目录提供 `main.go` 串联主题子包，主题子包通过导出函数返回可断言结果，测试采用表驱动 + `t.Run`，README 与 EXERCISES 补充学习说明。

第 04 章需要覆盖的知识点跨度较大：接口、错误处理、泛型。为避免一个文件堆叠过多概念，设计上沿用第 03 章的多子包结构，并保持所有示例只依赖标准库。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 04 章入口程序，输出接口、错误、泛型主题报告。
- 通过独立子包组织 interface、error、generic 三类核心概念。
- 所有关键行为都能通过返回值和单元测试断言，而不是只能人工阅读输出。
- README 与练习题替换占位内容，形成完整学习单元。
- 保持 `go test ./...` 与 `go build ./...` 通过。

**Non-Goals:**

- 不引入第三方泛型约束库或实验包。
- 不深入并发、标准库工程化或业务错误码体系，这些留给后续章节。
- 不修改第 01~03 章已经落地的源码行为。
- 不在本章引入复杂项目结构或外部 I/O。

## Decisions

### 1. 使用三个主题子包承载核心概念

- `iface/`：演示隐式实现、小接口、`any`、类型断言、type switch。
- `apperr/`：演示 sentinel error、自定义错误类型、`fmt.Errorf("%w")` 包装、`errors.Is` 与 `errors.As`。
- `generic/`：演示类型参数、类型集约束、泛型 `Map` / `Filter` 等基础用法。

**Rationale:** 与第 03 章 `seq` / `dict` / `model` / `ptr` 的结构一致，便于学习者按主题阅读和测试。

**Alternative considered:** 全部放在根目录单个 `main.go`。该方案更短，但不利于测试，也无法体现包组织与接口设计原则。

### 2. 示例函数返回结构化结果，入口程序只负责展示

主题子包提供导出函数返回字符串、结构体或切片，测试直接断言这些返回值。`main.go` 只组合这些函数并打印报告。

**Rationale:** 学习示例既能运行，也能测试；避免测试依赖脆弱的完整 stdout 文本。

**Alternative considered:** 示例函数直接 `fmt.Println`。这更直观但更难覆盖边界场景。

### 3. 泛型约束仅使用标准库可表达的类型集

使用本包内定义的 `Number` / `OrderedText` 等约束，避免依赖 `golang.org/x/exp/constraints`。

**Rationale:** 当前仓库无第三方依赖，章节学习重点是泛型语法与约束思想，不需要引入依赖管理噪音。

**Alternative considered:** 使用 `cmp.Ordered` 或外部 constraints。考虑到 Go 版本与教学简洁性，本章优先自定义最小约束。

## Risks / Trade-offs

- [Risk] 示例过多导致初学者难以聚焦 → Mitigation：每个子包只覆盖本章关键知识点，并用 README 串联学习顺序。
- [Risk] 错误处理示例过于玩具化 → Mitigation：用查询用户这类简单领域展示 sentinel + wrap + custom type，保证 `Is` / `As` 的价值清晰。
- [Risk] 泛型约束涉及类型集语法，初学者可能陌生 → Mitigation：只实现 `Map` / `Filter` / `Sum` 这类最小函数，并在练习中逐步扩展。
- [Risk] 章节目录从占位变为已落地，需要同步学习课程规格 → Mitigation：在 `learning-curriculum` delta 中明确第 04 章已落地后可包含源码、测试与练习。
