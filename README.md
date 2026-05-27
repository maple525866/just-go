# just-go · 从小白到 Go 架构师的学习仓库

> 一份"教科书式"的 Go 学习路线，覆盖 **语法 → 业务 → 架构** 三阶段共 **15 章 + 3 个综合项目**。

## 这是什么

`just-go` 不是一个产品仓库，而是**一份可执行的 Go 学习课表**。

- 每一章是一个独立可运行的 demo（按主题分目录）。
- 每一阶段以一个 capstone 综合项目收尾，把所学知识点串成业务。
- 借助 [OpenSpec](https://github.com/openspec-dev/openspec) 工作流，**每一章 = 一个独立的 OpenSpec change**，从而保证学习路径的可追溯、可回滚、可复盘。

## 适合谁

- ✅ **L0 完全零基础**：装好 Go 即可从第 01 章入门。
- ✅ **L2 略懂语法**：可跳过前 3 章，从并发 / 标准库开始。
- ✅ **L3 业务薄弱**：直接进入阶段二的 Web / 数据库 / 可观测性。
- ✅ **L4 缺架构视角**：直接进入阶段三的整洁架构 / DDD / 微服务。

## 快速开始

1. 阅读 [`ROADMAP.md`](./ROADMAP.md) —— 这里有完整的 15 章路线 + 3 个 capstone 项目说明。
2. 根据自己的水平定位起点章节。
3. 想正式开学某一章时，运行：
   ```
   /opsx-propose chapter-NN-<kebab-name>
   ```
   例如 `/opsx-propose chapter-01-hello-go` 会自动起一个 OpenSpec change 来落地第 01 章。

## 学习方式

- **章节式（推荐）**：从 01 开始按序推进，每完成一章勾选 [`ROADMAP.md`](./ROADMAP.md#进度追踪) 中对应 checkbox。
- **跳学**：根据 ROADMAP 中每章的"🔗 前置依赖"自行规划路径。

## 目录速览

```text
just-go/
├── ROADMAP.md                ★ 完整学习路线图
├── stage-1-syntax/           阶段一：Go 语法精通（7 章 + capstone-1）
├── stage-2-business/         阶段二：Go 业务工程（4 章 + capstone-2）
├── stage-3-architecture/     阶段三：Go 架构进阶（4 章 + capstone-3）
├── docs/                     跨章节文档（术语表 / FAQ / 参考资料）
└── openspec/                 OpenSpec 变更管理（每章一个 change）
```

## 进度追踪

详见 [`ROADMAP.md` 的进度追踪表](./ROADMAP.md#进度追踪)。

## License

[MIT](./LICENSE)
