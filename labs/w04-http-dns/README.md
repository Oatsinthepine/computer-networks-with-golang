# W4 — HTTP and DNS

## 实验

```bash
dig homepage.ntu.edu.tw
curl -v http://example.com/
curl -I https://example.com/
```

记录 DNS answer、目标 IP、TCP/TLS 建连提示、HTTP method、status 和关键 headers。不要把 cookies、authorization 或其他敏感 header 粘贴到仓库。

## Cloud 对照

说明浏览器访问 Ingress 域名时，CoreDNS、外部 DNS、Load Balancer 和 HTTP router 分别位于哪一段。
