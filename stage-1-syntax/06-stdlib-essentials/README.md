# 06. 标准库精要

> 阶段：① 语法精通 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：2 天

## 🎯 学习目标

熟练使用 Go 标准库中最常用的 7 个包。

## 🧩 关键知识点

- `fmt`：格式化输入输出
- `io` / `bufio`：读写抽象与缓冲
- `os` / `os/exec`：文件与进程
- `net` / `net/http`：网络与 HTTP（仅基础，深入留给阶段二）
- `encoding/json` / `encoding/xml`：序列化
- `time`：时间、定时器、ticker
- `reflect`：反射（只学读，不滥用）

## 📦 本章产出

本章提供一组可运行、可测试的标准库精要示例：

```text
stage-1-syntax/06-stdlib-essentials/
├── main.go               # 组装标准库精要学习报告
├── main_test.go          # 入口报告测试
├── format/               # fmt.Sprintf / fmt.Fprintf
├── stream/               # io.Reader / io.Writer / io.Copy / bufio.Scanner
├── system/               # os 临时文件 / 环境变量 / os/exec
├── web/                  # net/http handler / client / httptest
├── codec/                # encoding/json / encoding/xml round trip
├── clock/                # time.Format / Duration / Ticker
├── inspect/              # reflect 只读类型、字段、tag 检查
└── EXERCISES.md          # 课后练习与验收标准
```

运行示例：

```bash
go run ./stage-1-syntax/06-stdlib-essentials
```

运行本章测试：

```bash
go test ./stage-1-syntax/06-stdlib-essentials/...
```

## 🔗 前置依赖

- 第 05 章

## 📚 推荐扩展阅读

- [pkg.go.dev/std](https://pkg.go.dev/std)
- 《Go 程序设计语言》第 5、10、11、12 章

## ✅ 自测清单

- [ ] 能用 `fmt.Sprintf` 与 `fmt.Fprintf` 生成格式化输出
- [ ] 能用 `io.Reader` / `io.Writer` 写一个流式拷贝程序
- [ ] 能用 `bufio.Scanner` 按行读取文本
- [ ] 能用 `os` 安全读写临时文件，并读取环境变量
- [ ] 能用 `os/exec` 执行外部命令并处理输出 / 错误
- [ ] 能用 `net/http` 编写 handler，并用 `httptest` 做本地测试
- [ ] 能用 `encoding/json` 与 `encoding/xml` 完成结构体 round trip
- [ ] 能用 `time.Format`、`time.Duration`、`time.Ticker` 处理时间场景
- [ ] 能用 `reflect` 只读检查结构体字段与 tag，并说出反射的禁用场景
