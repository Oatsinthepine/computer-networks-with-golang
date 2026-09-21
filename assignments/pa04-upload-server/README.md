# PA4 — Single-upload server

本目录实现 PA4：一个只处理一次上传的 TCP 服务端。它先读取 PA3 client 发送的文件大小行，再逐行读取正文，在每行前加行号并写入 `whatever.txt`，最后回复原文件和新文件的字节大小。

## 运行

从本目录运行，以便 `whatever.txt` 生成在这里：

```sh
cd assignments/pa04-upload-server
go run PA4.go
```

再在项目根目录启动 PA3 client：

```sh
go run ./assignments/pa03-upload-client/PA3.go
```

输入：

```text
./assignments/pa03-upload-client/testdata/upload.txt
```

PA4 监听 `127.0.0.1:12000`，与当前 PA3 client 默认地址一致。每次上传后 PA4 会退出，下一次测试须重启服务端。`whatever.txt` 会被覆盖，只用于练习。

## 学习目标

- TCP listener、accept 与 connection 生命周期。
- 从 socket 读取协议头和文件内容。
- 文件写入与 client/server 协作。

## 检查清单

- [x] 完成 `examples/` 中的 Reader.ReadString 讲义示例。
- [x] 创建 `PA4.go`。
- [x] 使用 PA3 client 做本机端到端验证。
- [x] 将验证结果写入 `evidence/README.md`。
