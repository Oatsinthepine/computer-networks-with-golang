# PA3 本地上传接收端

这是给 PA3 客户端配套的小文件练习服务器，不替代 PA4 作业。它接收一次上传，打印内容并回复，然后退出；不会保存或覆盖文件。仅监听本机，文件大小限制为 1 MiB，没有超时设置，等待时可用 Ctrl+C 退出。

在项目根目录打开两个终端。

终端 A 先运行：

```sh
go run ./assignments/pa03-upload-client/examples/upload_server
```

终端 B 再运行：

```sh
go run ./assignments/pa03-upload-client
```

出现 `Upload filename:` 后输入：

```text
./assignments/pa03-upload-client/testdata/upload.txt
```

每次重新测试都需要先重启服务器。若提示端口已占用，检查是否已经运行了一个服务器，不要重复启动。

协议：一行十进制字节数 + 对应数量的原始字节。读取大小行和正文共用同一个 bufio.Reader，避免遗漏缓冲区中的正文；io.ReadFull 不要求一次底层读取就收到所有数据。
