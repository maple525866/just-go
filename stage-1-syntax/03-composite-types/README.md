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

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-03-composite-types` 来落地本章内容。

## 🔗 前置依赖

- 第 02 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 4 章
- [Go Slices: usage and internals](https://go.dev/blog/slices-intro)

## ✅ 自测清单（落地后填充）

- [ ] 能默写 `slice` 扩容的策略
- [ ] 能解释为什么 map 不能安全并发写
- [ ] 能讲清值接收者 vs 指针接收者的选择标准
- [ ] 能演示一次"切片共享底层数组"踩坑
- [ ] 能用 struct 嵌入实现"组合优于继承"
