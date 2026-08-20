# W9 — TCP observation

运行 PA5 或 PA6 后，使用系统工具观察监听和连接：

```bash
lsof -nP -iTCP -sTCP:LISTEN
netstat -an | grep tcp
```

记录 LISTEN、ESTABLISHED、TIME_WAIT 的含义，并解释为什么 application 的一次“上传”可能涉及多次 TCP 状态变化。
