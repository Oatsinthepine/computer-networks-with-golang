package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// 这个练习的服务端目前是只接受一次客户端上传文件后就退出，因为目前main 中只有一次 listener.Accept()， 没有循环，处理完回复 main 就结束了， defer 会关闭连接和监听
// 因此之后的 PA5 才会把接受连接的过程放进循环，PA6再引入并行处理。
func main() {
	// 1. 监听本机端口，等待一个客户端连接。
	listener, err := net.Listen("tcp", "127.0.0.1:12000")
	check(err)
	defer listener.Close()
	fmt.Println("Waiting for an upload on 127.0.0.1:12000...")

	// Accept() 等待一个客户端连接, 当客户端运行Dial("tcp", "127.0.0.1:12000")后，服务端 Accept()返回一个conn
	conn, err := listener.Accept()
	check(err)
	defer conn.Close()

	// 2. 读取上传的文件
	reader := bufio.NewReader(conn)
	sizeLine, err := reader.ReadString('\n')
	check(err)
	size, err := strconv.Atoi(strings.TrimSpace(sizeLine))
	check(err)
	// 教学版一次把内容放入内存，限制为 1 MiB 的小文件。
	if size < 0 || size > 1024*1024 {
		fmt.Println("Invalid file size: expected 0 to 1048576 bytes.")
		return
	}

	// 3. TCP 没有消息边界，ReadFull 会读满指定字节数或返回错误。 必须继续使用同一个 reader：它可能已经缓存了正文的一部分。
	data := make([]byte, size)
	// 真正的文件接受在这里
	_, err = io.ReadFull(reader, data)
	check(err)
	fmt.Printf("Received %d bytes. File content:\n%s\n", len(data), data)

	// 4. 发回一行应用层回复，不将文件写入磁盘。
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(fmt.Sprintf("%d bytes received\n", len(data)))
	check(err)
	check(writer.Flush())
	// 只处理一次上传；main 返回时关闭连接和监听端口。
}
