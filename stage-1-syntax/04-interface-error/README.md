# 04. 接口与错误

> 阶段：① 语法精通 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：1 天

## 🎯 学习目标

理解 interface 隐式实现，掌握标准错误处理，初识泛型。

## 🧩 关键知识点

- interface 隐式实现、空接口 `any`、类型断言、类型 switch
- 接口设计原则（小接口、接受接口返回结构体）
- `error` 接口、`errors.New` / `fmt.Errorf` / `%w` 包装
- `errors.Is` / `errors.As` 错误判别
- 泛型基础：类型参数、约束（constraints）、基本用法

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-04-interface-error` 来落地本章内容。

## 🔗 前置依赖

- 第 03 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 7 章
- [Go 1.18 泛型博客](https://go.dev/blog/intro-generics)
- [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)

## ✅ 自测清单（落地后填充）

- [ ] 能解释为什么 Go 接口是"隐式实现"
- [ ] 能讲清 `errors.Is` 与 `errors.As` 的差别
- [ ] 能写一个带 `%w` 包装的自定义错误链
- [ ] 能用泛型重构一个 `slice.Map` / `slice.Filter`
- [ ] 能说出 `any` vs `interface{}` 的关系
