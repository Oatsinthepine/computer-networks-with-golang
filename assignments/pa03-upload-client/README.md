# PA3 — TCP upload client

本目录包含 PA3 TCP 上传客户端和本地协议测试。实现顺序为：连接服务器、输入文件名、发送十进制文件大小与换行、发送原始文件字节、读取并打印回复、关闭连接。

## 运行

在项目根目录执行：

```sh
go run ./assignments/pa03-upload-client
go test -race ./assignments/pa03-upload-client
```

运行客户端前需要启动匹配上传协议的接收端。现有 simple_server 只接收一行，不是文件上传服务器；自动化测试会自行启动临时回环接收端，无须学校账号或外部服务。

提示文件名时可输入 `assignments/pa03-upload-client/testdata/upload.txt`，或完整路径（含空格时不需要额外引号）。相对路径以运行时工作目录为准。

地址在 PA3.go 的 serverAddress 中设置，默认 `127.0.0.1:12000`。教学版没有命令行参数或超时设置；服务器不回复时可能一直等待，可用 Ctrl+C 退出。文件一次性读入内存，仅适用于小文件练习。回复按行读取，使用 Scanner 默认长度限制，不用于超长回复。未验证学校服务器兼容性。

## 阅读代码

- 从 main 自上而下阅读，check(err) 与讲义一样，出错就停止。
- `os.ReadFile` 读取文件为 []byte，`len(data)` 获取字节大小。
- `fmt.Sprintf` 生成大小行，WriteString 写大小行，Write 写原始字节，不添加行号或正文换行。
- 大小行与正文使用同一个缓冲 writer，最后 `Flush()`。
- 服务器回复是应用层消息，不是 TCP ACK。
- `PA3_test.go` 的接收端只用于协议测试，不代替 PA4 作业。

## 验收

- [ ] 完成讲义中的简单 client/server 示例。
- [ ] 独立创建 `PA3.go`。
- [ ] 使用本机回环地址完成上传验证。
- [ ] 不向仓库提交真实或敏感上传内容。
