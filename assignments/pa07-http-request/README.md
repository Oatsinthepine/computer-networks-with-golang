# PA7 — HTTP request interpreter

已实现 native/PA7.go 与 nethttp/PA7.go。学习状态：实现与自动化验证完成，个人手动验收待记录。

先读原生版本，再对比标准库版本。两版均循环处理连接、仅支持 GET HTTP/1.x、在终端打印文件字节数或 File not found，随后关闭连接，不生成 HTTP response。

分别在 native 或 nethttp 目录运行 go run PA7.go，监听 127.0.0.1:12000；同一时间只启动一个版本。两版共用 native/testdata 样例。

测试：go test -race ./assignments/pa07-http-request/native ./assignments/pa07-http-request/nethttp（从项目根目录）。
