# 16 单元路线

完成一项后把 `[ ]` 改为 `[x]`。日期由你实际开始时填写，不绑定 NTU 正式学期。

| 状态 | 单元 | 理论重点 | 实践产出 |
|------|---:|---|---|
| [x]  | W1 | Internet、protocol、分层、edge/core | 五层网络地图和 GoLand 检查 |
| [x]  | W2 | Packet switching、delay、loss、throughput | PA1 waived（已有 Unix 基础） |
| [x]  | W3 | Application principles、process、QoS | PA2 文件 I/O |
| [x]  | W4 | HTTP、DNS、request/response | curl/DNS micro-lab |
| [x]  | W5 | SMTP、P2P、Multimedia 基础 | PA3 TCP 上传客户端 |
| [x]  | W6 | Application Layer 复盘 | 重测 PA2–PA3 |
| [x]  | W7 | Transport、UDP、stop-and-wait | PA4 单次上传服务器 |
| [x]  | W8 | Pipelining、GBN、SR | PA5 循环服务器 |
| [ ]  | W9 | TCP 重传、流控、连接管理 | PA6 并发服务器 |
| [ ]  | W10 | Congestion control | PA7 HTTP 请求解析 |
| [ ]  | W11 | Transport 复盘 | 并发、断连和错误输入测试 |
| [ ]  | W12 | Router、IPv4、DHCP、CIDR | 子网与路由 micro-lab |
| [ ]  | W13 | NAT、IPv6、SDN data plane | PA8 HTTP 文件服务器 |
| [ ]  | W14 | LS、DV、routing principles | PA9 TLS Web server |
| [ ]  | W15 | OSPF、BGP、SDN control plane、Ethernet | 端到端排障流程 |
| [ ]  | W16 | 全课程整合 | 全量测试和个人总结 |

## PA 进度记录

- PA1：waived。已有 Unix、终端和远程环境使用基础，未重复完成 PA1 指定练习。
- PA2：2026-08-22 completed（learning-goal standard）。讲义示例与最终程序已完成并验证.
- PA3: 2026-9-20 completed。
- PA4: 2026-9-21 completed。
- PA5: 2026-9-22 completed。

## 三个阶段检查点

### W6：Application

- 能解释 socket 为什么也是一种 I/O。
- 能从浏览器动作追踪到 DNS、TCP 和 HTTP。
- PA2、PA3 的自动化测试通过。

### W11：Transport

- 能解释可靠性、流量控制和拥塞控制的区别。
- 能说明顺序服务器与并发服务器的等待差异。
- PA4–PA7 通过普通测试和 race detector。

### W16：Network and integration

- 能解释 IP forwarding 与 routing 的区别。
- 能从 DNS → TCP → TLS → HTTP → application 定位访问失败环节。
- 全项目 `go test ./...` 和 `go test -race ./...` 通过。
