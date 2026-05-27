# 学习常见问题（FAQ）

> 本文件由学习过程中的真实疑问增量积累。每个章节 change 在落地时 SHALL 把本章的真实坑点追加一条到此处（按章节聚类）。

---

## 通用环境与工具

### Q1. 为什么 `go run main.go` 提示 "package main is not in std" 或找不到包？

**A：** 大多数情况是因为：

1. 当前目录不在某个 module 内（缺少 `go.mod`）。解决：在项目根目录执行 `go mod init <module-name>`。
2. 你写的 `import` 路径与 `go.mod` 的 `module` 声明不一致。请把 import 路径前缀替换为 `go.mod` 中 `module` 后面那个路径。
3. 不要在 GOPATH 下工作（Go 1.16+ 默认 module 模式）。

---

<!--
追加模板（请在合适的聚类下方追加）：

### QN. 一句话提问？

**A：** 一段话答案 + 可选代码示例。
-->
