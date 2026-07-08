# 08. Web 基础

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：2 天

## 🎯 学习目标

用 `net/http` 起服务，掌握路由、中间件、参数校验、JSON 编解码。

## 🧩 关键知识点

- `net/http` Server / Handler / ServeMux
- 主流 router（本章使用 chi）的方法路由、URL 参数与中间件链
- 中间件链：Recover、请求日志、CORS、简易限流、request ID/context 传播
- 参数校验：`go-playground/validator`
- JSON 编解码与 RESTful 状态码约定

## 📦 本章产出

本章实现了一个内存版博客文章 HTTP API，展示从标准库到 chi router 的 Web 服务基础：

```text
stage-2-business/08-web-foundations/
├── main.go                  # 启动 HTTP 服务，ADDR 未设置时默认 :8080
├── model/                   # Article、CreateArticleRequest、ErrorResponse 等 JSON 类型
├── store/                   # 内存文章仓库，提供确定性 seed 数据
├── response/                # JSON 成功响应与 REST 错误响应封装
├── validation/              # go-playground/validator 请求体验证封装
├── middleware/              # Recover / Logger / CORS / Limiter / RequestID 中间件
└── server/                  # 标准库 ServeMux 示例与 chi REST router
```

运行测试：

```bash
go test ./stage-2-business/08-web-foundations/...
```

启动服务：

```bash
go run ./stage-2-business/08-web-foundations
# 或指定端口
ADDR=:8081 go run ./stage-2-business/08-web-foundations
```

示例请求：

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/api/articles
curl http://localhost:8080/api/articles/1
curl -i -X POST http://localhost:8080/api/articles \
  -H 'Content-Type: application/json' \
  -d '{"title":"Testing handlers","body":"httptest makes handlers easy to verify.","tags":["test","http"]}'
```

## 🔗 前置依赖

- 第 07 章（阶段一全部）

## 📚 推荐扩展阅读

- [chi](https://github.com/go-chi/chi)
- [gin](https://gin-gonic.com/)
- [Building Web Apps with Go](https://gopherguides.com/articles)
- [go-playground/validator](https://github.com/go-playground/validator)

## ✅ 自测清单

- [ ] 能用标准库 `net/http` 和 `http.ServeMux` 写出 `GET /healthz`。
- [ ] 能解释 chi router 如何通过 `/{id}` URL 参数定位资源。
- [ ] 能写一条标准中间件链，并说明 Recover / 日志 / CORS / 限流 / request ID 各自的位置。
- [ ] 能用 validator tag 校验 JSON 请求体，并把校验错误转成 422 响应。
- [ ] 能正确处理 `context.Context` 在 HTTP 请求中的传播。
- [ ] 能说清 `400`、`404`、`422`、`429`、`500` 在本章 REST API 中的使用场景。
