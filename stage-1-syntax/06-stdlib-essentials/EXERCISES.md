# 第 06 章练习：标准库精要

> 建议先运行 `go test ./stage-1-syntax/06-stdlib-essentials/...`，确认示例代码全部通过，再开始练习。

## 练习 1：实现一个带统计信息的流式拷贝

在 `stream` 包中新增函数，接收 `io.Reader` 和 `io.Writer`，使用 `io.Copy` 完成复制，并返回复制字节数与目标内容摘要。

**验收标准：**

- 使用 `strings.Reader` 与 `bytes.Buffer` 编写测试，不读写真实固定文件。
- 输入 `hello stdlib` 时复制字节数与输出内容均可断言。
- 使用表驱动 + `t.Run` 编写测试。
- 运行 `go test ./stage-1-syntax/06-stdlib-essentials/stream` 通过。

## 练习 2：扩展 HTTP JSON handler

在 `web` 包中新增一个 handler，返回 `Content-Type: application/json` 和一段 JSON 响应，例如 `{"message":"ok"}`。

**验收标准：**

- 测试必须使用 `httptest.NewRecorder` 或 `httptest.NewServer`。
- 断言状态码、Content-Type 与 JSON 响应体。
- 不依赖外部互联网服务。
- 运行 `go test ./stage-1-syntax/06-stdlib-essentials/web` 通过。

## 练习 3：实现自定义 JSON 字段校验

在 `codec` 包中新增一个函数，接收 JSON 字符串并解码为结构体；当必填字段为空时返回错误。

**验收标准：**

- 合法 JSON 能成功解码。
- 缺少标题或标题为空时返回可断言错误。
- 非法 JSON 返回解码错误。
- 使用表驱动 + `t.Run` 编写测试。

## 练习 4：实现有限心跳器

在 `clock` 包中新增一个函数，使用 `time.Ticker` 产生指定次数的心跳文本，例如 `tick-1`、`tick-2`。

**验收标准：**

- `n=3` 时返回 3 个心跳文本。
- `n=0` 时立即返回空切片。
- 测试不依赖精确耗时，只断言返回数量和内容。
- 运行 `go test ./stage-1-syntax/06-stdlib-essentials/clock` 通过。

## 练习 5：读取更多 struct tag

在 `inspect` 包中扩展字段描述，让它同时读取 `json` 与 `xml` tag。

**验收标准：**

- 示例结构体至少包含两个字段，且字段同时带 `json` / `xml` tag。
- 测试断言字段名、类型、json tag、xml tag。
- 不使用反射修改字段值。
- 运行 `go test ./stage-1-syntax/06-stdlib-essentials/inspect` 通过。
