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

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-06-stdlib-essentials` 来落地本章内容。

## 🔗 前置依赖

- 第 05 章

## 📚 推荐扩展阅读

- [pkg.go.dev/std](https://pkg.go.dev/std)
- 《Go 程序设计语言》第 5、10、11、12 章

## ✅ 自测清单（落地后填充）

- [ ] 能用 `io.Reader` / `io.Writer` 写一个流式拷贝程序
- [ ] 能用 `encoding/json` 实现自定义 Marshal / Unmarshal
- [ ] 能用 `time.Ticker` 实现一个心跳器
- [ ] 能用 `os/exec` 执行外部命令并管道传递输出
- [ ] 能说出 `reflect` 的合理用例与禁用场景
