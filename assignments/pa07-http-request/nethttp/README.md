# net/http 对照

从本目录运行 go run PA7.go，默认网站目录为 ../native/testdata。
先停止原生版，两版使用同一个 127.0.0.1:12000。测试命令见 ../native/START-HERE.md。

main 与文件查找逻辑保持独立可读，不抽取共享网络库。
serveRequest 用 http.ReadRequest(reader) 替代原生版手动拆解请求行和 headers：
req.Method 对应 method，req.RequestURI 对应 request-target，req.Host 对应 Host。
标准库负责 HTTP 语法解析；应用仍决定支持的 method、文件访问范围和日志。

仍使用 net.Listen/Accept，只打印文件信息并关闭连接。没有使用 http.ListenAndServe，
因此不会由 HTTP server 自动生成响应；curl 的空回复仍是预期的 PA7 行为。
