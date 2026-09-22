# PA5 — Looping server

状态：练习中。PA5.go 和 client/main.go 是可编译的 TODO 骨架，
目前运行时只打印未完成提示后退出，不会监听、上传或生成输出。

## 练习顺序

1. 完成 client/main.go：参考 PA3，连接后输入文件名，发送大小行和正文，刷新，再读回复。
2. 完成 PA5.go 的 handleUpload(conn) error：参考 PA4，逐行编号并写入 whatever.txt，回复原始与新文件大小。
3. 完成服务端 main：循环外 Listen，循环内 Accept 并同步调用 handleUpload。单次上传失败打印错误，继续处理下一次。

两端约定使用 127.0.0.1:12000。协议仍为“十进制大小行 + 原始正文”，
回复为一行文本；不使用 HTTP、Gin 或 goroutine。仅接收不超过 1 MiB 的小文件。
每次上传重置行号为 1，并覆盖当前运行目录的 whatever.txt。

## 完成 TODO 后运行

两个终端最初都位于项目根目录。先确认 PA4 或 PA3 测试服务器没有占用 12000。

终端 A：

```sh
cd assignments/pa05-looping-server
go run PA5.go
```

终端 B：

```sh
cd assignments/pa05-looping-server
go run ./client
```

客户端提示文件名时输入 `testdata/input.txt`。输出将出现在
服务端运行目录的 whatever.txt；该文件会覆盖旧内容，已由仓库规则忽略。
client 使用独立目录，避免与 PA5.go 的 main 函数冲突。

## 连接生命周期

listener 在服务端存活期间持续存在，每次 conn 在 handleUpload 返回时关闭。
你已完成的 examples/loop_server/server-loop.go 保留作为对照：
其中循环内的 defer conn.Close() 会等到整个 main 返回才执行，并非每轮结束执行。
在 handleUpload 内使用 defer，便能让每次上传结束时及时释放资源。

PA5 是顺序处理。第二个客户端可能已完成 TCP 连接并进入系统等待队列，
但服务端应用仍需要先处理完第一个上传。连接建立不等于文件已被处理。

## 后续验收（当前均未执行）

- [ ] 连续启动客户端三次，服务器无需重启；每次输出重新编号，不追加旧内容。
- [ ] 逐字节核对原始大小、编号后的大小及回复。
- [ ] 空文件输出为空；空行也编号；中文按字节计数；末行无换行仍能完成。
- [ ] 客户端提前断开后，下一次正常上传仍成功。
- [ ] 客户端 A 连接后停在文件名输入处；启动 B 并输入文件名，观察 B 等待。
      完成 A 的输入上传后，确认 B 随后得到回复。

这里只做骨架编译检查。待你完成 TODO 后，再记录上述端到端结果。

## 检查清单

- [ ] 复制并理解自己的 PA4，而不是使用预制答案。
- [ ] 自己补完 `PA5.go` 和独立客户端的 TODO。
- [ ] 连续运行多个 client，server 无需重启。
- [ ] 记录顺序处理带来的等待现象。
