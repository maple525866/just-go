# 08. Web 基础练习

## 练习 1：新增一个标准库 ServeMux 路由

在 `server.NewStdMux` 中新增 `GET /version`，返回：

```json
{"version":"chapter-08"}
```

**验收标准：**

- 新增 `httptest` 测试覆盖 `GET /version`。
- 执行 `go test ./stage-2-business/08-web-foundations/server` 通过。
- 响应头 `Content-Type` 为 `application/json`。

## 练习 2：为文章列表增加 tag 过滤

让 `GET /api/articles?tag=http` 只返回包含 `http` 标签的文章。

**验收标准：**

- 不传 `tag` 时仍返回全部文章。
- 传入已有 tag 时只返回匹配文章。
- 传入不存在的 tag 时返回空列表且状态码仍为 200。
- 使用表驱动测试覆盖以上三种情况。

## 练习 3：扩展 validator 校验规则

为 `CreateArticleRequest.Title` 增加最小长度 3 的限制。

**验收标准：**

- `title` 长度小于 3 时，`POST /api/articles` 返回 422。
- 错误响应中包含字段名 `title` 和规则名 `min`。
- 原有合法创建文章测试仍通过。

## 练习 4：为中间件链加入响应时间头

新增一个 middleware，在响应头写入 `X-Response-Time`，值可以是类似 `1ms` / `250µs` 的耗时字符串。

**验收标准：**

- middleware 类型为 `func(http.Handler) http.Handler`。
- 使用 `httptest` 验证响应头存在且非空。
- 将该 middleware 接入 `server.NewRouter` 的中间件链。

## 练习 5：补充 REST 错误场景说明

在 README 的自测清单下新增一张表，说明本章 API 中 `400`、`404`、`422`、`429`、`500` 分别由什么场景触发。

**验收标准：**

- 表格至少包含“状态码 / 场景 / 示例接口”三列。
- 每个状态码都能在现有测试或新增测试中找到对应断言。
