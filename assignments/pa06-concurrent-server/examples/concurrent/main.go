package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const serverAddress = "127.0.0.1:12000"
const processingDelay = 5 * time.Second

func main() {
	fmt.Println("PA6 并发延迟练习。")
	// TODO 1：循环外 Listen；检查错误，安排 listener.Close()。
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Printf("监听失败：%v\n", err)
		return
	}
	defer listener.Close()

	// TODO 2：循环内 Accept；失败时记录原因并退出。
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败：%v\n", err)
			break
		}
		// TODO 3：为每个 conn 启动 goroutine，显式传入 conn。
		// 在 goroutine 内检查 handleConnection 返回的 error。
		// 不要在 Accept 循环中 defer conn.Close()。
		go func() {
			if err := handleConnection(conn); err != nil {
				fmt.Printf("处理连接失败：%v\n", err)
			}
		}()
	}
}

func handleConnection(conn net.Conn) error {
	// TODO 4：defer conn.Close()；记录带时间、客户端地址的开始与结束日志。
	defer conn.Close()
	fmt.Printf("连接开始：%s\n", conn.RemoteAddr())
	// TODO 5：用 bufio.Reader.ReadString('\n') 读取一行。
	// 返回内容包含换行；打印消息，出错返回 error。
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("读取消息失败：%v", err)
	}
	fmt.Printf("收到消息：%s", message)

	// TODO 6：用 bufio.Writer 回复并检查 Flush。
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("Hello World!\n")
	if err != nil {
		return fmt.Errorf("回复消息失败：%v", err)
	}

	if writerFlushErr := writer.Flush(); writerFlushErr != nil {
		return fmt.Errorf("刷新回复失败：%v", writerFlushErr)
	}

	fmt.Println("已回复成功。")

	// TODO 7：在成功 Flush 后 Sleep(processingDelay)，然后再返回。
	time.Sleep(processingDelay)
	fmt.Println("处理完成")
	return nil
}
