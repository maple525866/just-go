# 12. 整洁架构

> 阶段：③ 架构进阶 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

本章用“创建并发布文章”这一条业务切片演示整洁架构。重点不是目录命名，而是依赖方向：稳定的业务规则不认识 HTTP、内存仓储或依赖注入工具。

## 🎯 学习目标

- 区分 `domain`、`usecase`、`interface`、`infrastructure` 四层职责
- 用端口（Go interface）让外层适配器依赖业务核心
- 使用构造函数和 Wire 在组合根完成编译期依赖注入
- 使用 mock 隔离 usecase 测试，并用测试守护架构边界

## 🧩 关键知识点

- 整洁架构与六边形架构的依赖规则
- 实体、应用用例、输入输出适配器与基础设施适配器
- 依赖倒置、构造函数注入与 Wire provider
- mock 测试、compare-and-swap 更新与架构边界测试

## 🧭 依赖方向

```text
HTTP 请求 -> interface/httpapi -> usecase -> domain
                                  ^
                                  |
                     infrastructure/memory

main / Wire 只负责把具体实现组装起来
```

依赖只能指向更内层。`domain` 不导入本章其他层；`usecase` 只依赖 `domain`；外层通过实现 `usecase.ArticleRepository` 接入。

## 📦 本章产出

一个可运行的文章发布服务，包含显式业务规则、应用端口、内存与 HTTP 适配器、Wire 组合根，以及覆盖各层和依赖方向的自动化测试。

### 包职责

| 目录 | 职责 |
| --- | --- |
| `domain/` | Article 实体、校验和发布状态转换 |
| `usecase/` | Repository 端口与创建、查询、发布流程 |
| `interface/httpapi/` | JSON/HTTP 到应用输入输出的转换 |
| `infrastructure/memory/` | 并发安全的内存仓储适配器 |
| `wire.go` / `wire_gen.go` | Provider 声明与已提交的生成结果 |
| `main.go` | 最外层启动入口 |

## ▶️ 运行与验证

```bash
go run ./stage-3-architecture/12-clean-architecture
go test ./stage-3-architecture/12-clean-architecture/...
```

服务监听 `:8080`：

```bash
curl -X POST http://localhost:8080/articles -H "Content-Type: application/json" -d '{"title":"Ports","body":"Adapters"}'
curl http://localhost:8080/articles/article-1
curl -X POST http://localhost:8080/articles/article-1/publish
```

如需重新生成 Wire 输出：

```bash
go generate ./stage-3-architecture/12-clean-architecture
```

正常构建直接使用已提交的 `wire_gen.go`，不要求本地预装 Wire CLI。Wire 只是自动写构造函数调用；依赖倒置来自接口的位置和导入方向，不来自工具本身。

## 🧪 测试观察点

- `usecase/article_test.go` 使用手写 mock，不启动 HTTP 或仓储
- memory 测试验证并发安全仓储的副本隔离
- HTTP 测试验证协议转换与错误状态码
- `architecture_test.go` 解析 import，禁止核心层反向依赖外层

## 🔗 前置依赖

- Capstone 2：已有真实博客业务代码，便于比较重构前后的依赖方向

## 📚 推荐扩展阅读

- 《架构整洁之道》Robert C. Martin
- [Google Wire](https://github.com/google/wire)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)

## ✅ 自测清单

- [ ] 能画出本章四层依赖图
- [ ] 能解释为什么 domain 不能 import infrastructure
- [ ] 能说明 `ArticleRepository` 为什么放在 usecase 一侧
- [ ] 能在不启动数据库的情况下测试发布流程
- [ ] 能读懂 `wire.go` 与 `wire_gen.go` 的对应关系

练习见 [EXERCISES.md](./EXERCISES.md)。
