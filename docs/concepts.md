# 核心概念索引

## 分层地图

| 层 | 主要问题 | 课程关键词 | Cloud/Kubernetes 对照 |
|---|---|---|---|
| Application | 应用进程如何表达请求与响应 | HTTP、DNS、SMTP | API Gateway、CoreDNS、Ingress |
| Transport | 进程之间如何可靠或快速传输 | TCP、UDP、flow/congestion control | connection、timeout、gRPC |
| Network | 数据报如何跨网络转发 | IPv4/IPv6、NAT、routing | VPC、subnet、route table、CNI |
| Link | 相邻节点如何传递 frame | Ethernet、MAC、ARP | node/LAN 邻接网络 |
| Physical | bit 如何通过介质传递 | bandwidth、signal | 云平台通常隐藏此层 |

## 容易混淆的边界

- **Forwarding vs routing**：前者是单台路由器的本地转发动作，后者是决定端到端路径的网络级过程。
- **Flow control vs congestion control**：前者保护接收端，后者保护网络路径。
- **Concurrency vs capacity**：并发可降低客户端排队感受，但不自动增加 CPU、磁盘或链路总容量。
- **HTTP vs TCP**：HTTP 规定应用消息语义；TCP 提供有序、可靠的字节流。
- **TLS vs HTTP**：TLS 为传输中的应用数据提供身份验证、机密性和完整性；HTTPS 是 HTTP over TLS。

## 排障顺序

1. 名称能否解析：DNS/CoreDNS。
2. IP 与路由是否可达：subnet、route table、NAT。
3. TCP 是否建立连接：地址、端口、监听、firewall。
4. TLS 是否握手成功：证书、主机名、信任链、协议版本。
5. HTTP 是否得到预期状态：path、method、headers、status。
6. 应用是否健康：handler、文件、依赖服务和日志。
