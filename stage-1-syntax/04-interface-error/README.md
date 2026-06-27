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

## 📦 本章产出

本章提供一组可运行、可测试的接口、错误与泛型示例：

```text
stage-1-syntax/04-interface-error/
├── main.go                 # 组装接口、错误与泛型学习报告
├── main_test.go            # 入口报告测试
├── iface/
│   ├── iface.go            # 隐式实现、小接口、any、类型断言、type switch
│   └── iface_test.go       # iface 表驱动测试
├── apperr/
│   ├── apperr.go           # sentinel error、%w、errors.Is、errors.As
│   └── apperr_test.go      # apperr 表驱动测试
├── generic/
│   ├── generic.go          # 泛型 Map / Filter / Sum 与类型集约束
│   └── generic_test.go     # generic 表驱动测试
└── EXERCISES.md            # 课后练习与验收标准
```

运行示例：

```bash
go run ./stage-1-syntax/04-interface-error
```

运行本章测试：

```bash
go test ./stage-1-syntax/04-interface-error/...
```

## 🔗 前置依赖

- 第 03 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 7 章
- [Go 1.18 泛型博客](https://go.dev/blog/intro-generics)
- [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)

## ✅ 自测清单

- [ ] 能解释为什么 Go 接口是“隐式实现”，并指出 `iface.Book` 如何满足 `iface.Describer`
- [ ] 能说出“小接口”和“接受接口返回结构体”的设计收益
- [ ] 能讲清 `any` 与 `interface{}` 的关系，并读懂 type switch 分支
- [ ] 能写出带 `%w` 包装的错误链，并用 `errors.Is` 判断 sentinel error
- [ ] 能用 `errors.As` 从错误链中提取自定义错误类型
- [ ] 能用类型参数实现一个泛型 `Map` / `Filter`
- [ ] 能定义一个简单类型集约束，并解释 `generic.Sum` 为什么只能接收数值类型
