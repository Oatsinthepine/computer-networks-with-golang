# Gin file upload comparison lab

这个实验把 PA4 的“逐行编号”业务放到 HTTP API 中。它不会修改 PA3/PA4；目的是比较自定义 TCP 协议、HTTP 与 Gin 的抽象边界。

## 运行

在项目根目录执行：

```sh
go run ./labs/w04-http-dns/gin-file-upload
```

健康检查：

```sh
curl -i http://127.0.0.1:8080/health
```

上传现有测试文件：

```sh
curl -i -X POST http://127.0.0.1:8080/upload \
  -F "file=@./assignments/pa03-upload-client/testdata/upload.txt"
```

成功时返回 HTTP 200 和 JSON。处理结果写入运行目录下的 `outputs/whatever-*.txt`；文件名由服务器生成，避免信任客户端文件名和连续上传时互相覆盖。

## 错误响应

| 场景 | HTTP status |
|---|---:|
| 缺少 multipart 字段 `file` 或请求格式错误 | 400 |
| 文件超过 1 MiB | 413 |
| 无法创建输出目录、文件或转换失败 | 500 |
| 上传并转换成功 | 200 |

HTTP request 上限略高于 1 MiB，以容纳 multipart 边界和 headers；实际文件仍严格限制为 1 MiB。

## 与 PA4 的对应

| PA4 原生 TCP | Gin HTTP |
|---|---|
| `Listen`、`Accept` | `router.Run` 和 `net/http` 在内部管理 |
| 第一行传文件大小 | HTTP/multipart 描述 request body |
| `LimitedReader` 保护正文边界 | `MaxBytesReader` 和 `fileHeader.Size` 限制大小 |
| 从 `conn` 逐行读取 | 从上传文件的 `io.Reader` 逐行读取 |
| 写一行文本回复 | HTTP status + JSON body |
| 一次连接后退出 | HTTP server 持续处理请求 |

`addLineNumbers` 不依赖 Gin：Gin 负责 HTTP 输入输出，它负责 PA4 的业务转换。


