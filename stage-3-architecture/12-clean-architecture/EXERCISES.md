# Chapter 12 Exercises

## 1. 增加文章列表用例

为 repository 端口、memory adapter 和 HTTP adapter 增加文章列表。

验收标准：

- usecase 只依赖端口，不导入 infrastructure
- `GET /articles` 返回稳定顺序的 JSON 数组
- mock、memory、HTTP 三层测试均覆盖空列表与多文章场景

## 2. 替换持久化适配器

实现一个文件或 SQLite repository，并在组合根替换 memory provider。

验收标准：

- domain 和 usecase 无需修改
- 新适配器通过与 memory 相同的 repository 契约测试
- 重启进程后数据仍可读取，写入失败能返回可判断的错误

## 3. 增加文章归档规则

在 domain 增加归档状态转换，并通过 usecase 与 HTTP 暴露。

验收标准：

- 只有已发布文章可以归档，非法转换返回稳定 domain error
- usecase 测试只使用手写 mock
- HTTP 对非法转换返回 `409 Conflict`

## 4. 扩展架构守卫

让架构测试同时验证 infrastructure 不得导入 interface，并输出清楚的违规位置。

验收标准：

- 人为加入一个违规 import 时测试失败
- 删除违规 import 后测试恢复通过
- 测试不依赖第三方静态分析工具
