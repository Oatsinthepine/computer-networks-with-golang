# PA6 — Concurrent server

本目录用于从自己的 PA5 演进到 goroutine 并发 server。源码和测试有意留空。

## 两个练习分支

- 课程原版行为放在本目录的 `PA6.go`。
- 完成原版后，再在 `examples/strict/` 保留对照，并自行尝试避免并发文件覆盖的改进版。

## 检查清单

- [ ] 阅读本地 `PA6.pdf`。
- [ ] 独立提取连接处理函数。
- [ ] 独立添加 goroutine 和课程要求的延迟。
- [ ] 使用至少两个 client 比较 PA5 与 PA6。
- [ ] 自己添加 race detector 验证。
