# 12. 整洁架构

> 阶段：③ 架构进阶 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

## 🎯 学习目标

掌握分层（domain / usecase / interface / infrastructure）+ 依赖注入。

## 🧩 关键知识点

- 整洁架构 / 六边形架构核心思想
- 四层分层：domain / usecase / interface / infrastructure
- 依赖倒置原则：业务核心不依赖外部细节
- `wire` 编译期依赖注入
- 单元测试在分层架构下的玩法（mock 接口）

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-12-clean-architecture` 来落地本章内容。

## 🔗 前置依赖

- capstone-2（必须有"真业务代码"作为重构对象）

## 📚 推荐扩展阅读

- 《架构整洁之道》Robert C. Martin
- [wire](https://github.com/google/wire)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)

## ✅ 自测清单（落地后填充）

- [ ] 能画出本仓库的整洁架构四层依赖图
- [ ] 能解释为什么 domain 层不能 import infrastructure 层
- [ ] 能用 `wire` 写一份依赖注入配置
- [ ] 能用 mock 接口把 usecase 层的测试解耦于数据库
- [ ] 能讲清"依赖倒置"的语义与价值
