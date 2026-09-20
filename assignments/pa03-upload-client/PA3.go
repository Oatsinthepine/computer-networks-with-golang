package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var serverAddress = "127.0.0.1:12000"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 1. 连接服务器：主动建立 TCP 连接，对应服务器的 listener.Accept()
	// 连接成功后， conn 就是双向的通信的入口
	// 从 conn 读取 → 接收服务器发来的字节
	// 向 conn 写入 → 发送字节给服务器
	conn, err := net.Dial("tcp", serverAddress)
	check(err)
	defer conn.Close()

	// 2. 从键盘读取一整行文件名，路径可以包含空格。
	fmt.Print("Upload filename: ")
	// Os.Stdin 是控制台标准输入，来自键盘
	keyboard := bufio.NewScanner(os.Stdin)
	if !keyboard.Scan() {
		check(keyboard.Err())
		fmt.Println("No filename entered.")
		return
	}
	// .Test()读取控制台输入（字符串），此时只是获得文件路径，并没有打开文件。
	filename := keyboard.Text()

	// 3. 一次读取整个小文件；data 是 []byte，len(data) 是字节数。
	// 这里没有使用 Scanner 按行读取文件，因为我们练习模拟字节上传作为 packet
	data, err := os.ReadFile(filename)
	check(err)

	// 4. bufio.NewWriter(conn) 写向网络连接，这里就是连接的 127.0.0.12000的服务器
	writer := bufio.NewWriter(conn)
	sizeLine := fmt.Sprintf("%d\n", len(data))
	_, err = writer.WriteString(sizeLine)
	check(err)
	// 这里向服务器发送原始文件内容（字节）
	_, err = writer.Write(data)
	check(err)
	err = writer.Flush()
	check(err)

	// 5. 等待并显示服务器的回复
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		fmt.Println("Server replies:", scanner.Text())
	} else {
		check(scanner.Err())
		fmt.Println("Server closed without a reply.")
	}
	// main 结束时，defer 会关闭连接。
}
