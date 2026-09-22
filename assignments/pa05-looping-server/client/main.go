package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const serverAddress = "127.0.0.1:12000"

func main() {
	fmt.Println("PA5 客户端：上传单个小文件到服务器")

	// TODO 1: connect
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("连接服务器失败：", err)
		return
	}
	defer conn.Close()

	// TODO 2: read filename from stdin (whole line)
	fmt.Print("请输入要上传的文件路径：")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Println("读取文件名时出错：", err)
		} else {
			fmt.Println("没有读取到文件名，退出")
		}
		return
	}
	filename := strings.TrimSpace(scanner.Text())
	if filename == "" {
		fmt.Println("空文件名，退出")
		return
	}

	// TODO 3: read file
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}

	// TODO 4: write size + newline, then raw data
	writer := bufio.NewWriter(conn)
	_, err = fmt.Fprintf(writer, "%d\n", len(data))
	if err != nil {
		fmt.Println("写入长度失败：", err)
		return
	}

	n, err := writer.Write(data)
	if err != nil {
		fmt.Println("写入文件数据失败：", err)
		return
	}
	if n != len(data) {
		fmt.Printf("写入不完整：期望 %d 字节，写入 %d 字节\n", len(data), n)
		return
	}

	if err := writer.Flush(); err != nil {
		fmt.Println("刷新写缓冲失败：", err)
		return
	}

	// TODO 5: read one-line reply from server
	reader := bufio.NewReader(conn)
	reply, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取服务器回复失败：", err)
		return
	}

	fmt.Println("服务器回复：", strings.TrimSpace(reply))
}
