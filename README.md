# NTU Computer Networks with Go

这是台湾大学 901E31110 Introduction to Computer Networks (Spring 2026) 的低压自学工作区。目标不是追赶正式课程进度，而是通过 Go 练习建立 HTTP、TCP/UDP、IP、Routing 和 Ethernet 的完整心智模型。

## 使用方式

1. 用 GoLand 打开本目录；GoLand 会识别根目录的 `go.mod`。
2. 在 [`docs/roadmap.md`](docs/roadmap.md) 找到当前学习单元。
3. 阅读 `Lecture/` 中对应讲义，再完成 `labs/` 或 `assignments/` 中的实践。
4. 在 `docs/weekly/` 写下本周理解，在相应 `evidence/README.md` 记录人工验证。
5. 每个阶段运行：

   ```bash
   go test ./...
   go test -race ./...
   ```

## 目录边界

- `Lecture/`、`Go_network_assignment/`：原始课程资料，不修改。
- `docs/`：16 单元路线、概念索引和每周记录。
- `labs/`：不属于正式 PA 的短实验。
- `assignments/`：PA1–PA9 的空练习骨架；PA7–PA9 预留原生 socket 与 `net/http` 两条路线。
- `testdata/`：跨作业通用的无敏感测试输入。

## 运行约定

- 所有服务器默认只监听 `127.0.0.1`，避免意外暴露到局域网。
- 命令行参数通常按 `地址、文件目录或文件名` 的顺序传入；各 PA 的 README 有确切示例。
- TLS 证书只用于本机实验，不提交私钥。
- 当前项目不包含任何 PA 参考答案。源码和测试应在你开始对应单元时亲自创建。
- PA2–PA9 将保持自包含，以便以后看到每次作业带来的演进。

## 每周时间盒

- 30–45 分钟：讲义或视频。
- 15–20 分钟：写“解决什么问题、怎样解决”。
- 60–90 分钟：PA 或 micro-lab。
- 10 分钟：连接到 Cloud/Kubernetes 场景。

忙碌周可直接顺延。W6、W11、W16 是复盘与缓冲周，不形成“欠课”。
