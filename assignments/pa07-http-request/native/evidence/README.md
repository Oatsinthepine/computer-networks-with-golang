# PA7 native evidence

- 日期：2026-11-10
- agent 自动化检查：native 和 nethttp 的 go test -race 均通过。
- 存在文件：验证根路径映射、普通文件、子目录、查询参数及正确大小。
- 不存在文件：验证 File not found 日志。
- 畸形请求：验证无效请求行、无效 header、非 GET 拒绝。
- 文件范围：验证目录请求、编码的 .. 和目录外符号链接均拒绝。
- 协议：net.Pipe 测试验证请求完成后关闭连接，不发送 HTTP 响应。
- curl/browser：本次尚未进行真实 TCP curl 或浏览器验收，按 START-HERE.md 手动测试。
- 完成：实现及上述自动化检查完成；学习掌握和手动验收尚未标记完成。
