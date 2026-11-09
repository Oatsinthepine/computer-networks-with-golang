# PA7 运行与阅读

在项目根目录打开终端 A：

```sh
cd assignments/pa07-http-request/native
go run PA7.go
```

终端 B：

```sh
curl --http1.1 -v --max-time 5 http://127.0.0.1:12000/hello.html
curl --http1.1 -v --max-time 5 http://127.0.0.1:12000/nested/page.html
curl --http1.1 -v --max-time 5 http://127.0.0.1:12000/missing.html
curl --http1.1 -v --max-time 5 'http://127.0.0.1:12000/hello.html?source=lab'
wc -c assignments/pa07-http-request/native/testdata/hello.html
```

停止 PA6 等占用 12000 的服务器。PA5 上传客户端不发送 HTTP 请求，不能用作 PA7 客户端。

curl 的 Empty reply from server (52) 是预期现象：PA7 只解析并打印，不返回 HTTP 状态行或正文。超时不等于空回复。核对服务端日志及 wc -c 的字节数。

阅读顺序：
1. main：Listen、循环 Accept、同步调用 serveRequest。
2. serveRequest：读取请求行的三个字段，读 headers 到 CRLF 空行，检查 GET。
3. printFile：URL path 排除 query，/ 映射 hello.html，查找目录内的普通文件。
4. defer conn.Close：一条连接处理一次请求，之后继续接收下一条。

手工解析器仅用于学习 HTTP/1.x 子集，不是完整 HTTP 实现；不处理请求正文或 keep-alive。设有 10 秒读取期限和 16 KiB 请求头读取上限。

os.OpenRoot（Go 1.24+）把访问限制在 testdata 中，防止通过符号链接逃出目录；显式拒绝 .. 路径段。不存在打印 File not found，目录或畸形请求记录错误。未发送 200/404 响应，留待 PA8。
