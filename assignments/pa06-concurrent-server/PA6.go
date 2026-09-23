package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const serverAddress = "127.0.0.1:12000"
const maxFileSize = 1 << 20
const processingDelay = 5 * time.Second

func main() {
	fmt.Println("PA6 并发上传 Server launched.")
	// TODO 1：循环外 Listen，成功后 defer listener.Close()。
	// 监听或 Accept 失败时记录错误并退出。
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Printf("监听失败：%v\n", err)
		return
	}
	defer listener.Close()

	// TODO 2：循环内 Accept；显式传 conn 给一个新的 goroutine。
	// goroutine 内调用 handleUpload，检查并记录其返回错误。
	// 不要直接丢弃 handleUpload 的 error，也不要在循环里 defer conn.Close()。
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败：%v\n", err)
			break
		}
		go func(c net.Conn) {
			if err := handleUpload(c); err != nil {
				fmt.Printf("处理上传失败：%v\n", err)
			}
		}(conn)
	}
}

func handleUpload(conn net.Conn) error {
	// TODO 3：在此 defer conn.Close()，由本次处理负责连接生命周期。
	// 日志用时间和 conn.RemoteAddr() 标识“开始处理、回复已刷新、处理结束”。
	// 结束日志也应覆盖错误返回；日志交错是并发运行的正常现象。
	defer conn.Close()
	fmt.Printf("开始处理：%s\n", conn.RemoteAddr())
	// TODO 4：迁移 PA5 的大小行解析，验证 0 <= size <= maxFileSize。
	// 沿用同一个 reader，通过 LimitedReader 限定正文，再逐行处理。
	// reader、writer、行号、字节计数都在本次调用内创建。
	reader := bufio.NewReader(conn)
	sizeLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("读取文件大小失败：%w", err)
	}
	size, err := strconv.Atoi(strings.TrimSpace(sizeLine))
	if err != nil {
		return fmt.Errorf("解析文件大小失败：%w", err)
	}
	if size < 0 || size > maxFileSize {
		return fmt.Errorf("文件大小超出范围： %d", size)
	}

	// TODO 5：使用 os.CreateTemp(".", "whatever-*.txt") 为这次上传创建独立输出。
	// 记录 outputFile.Name()，让客户端地址与输出文件可对应。
	// 安排错误路径的文件关闭；成功路径显式检查 Close 的结果。
	outputFile, err := os.CreateTemp(".", "whatever-*.txt")
	if err != nil {
		return fmt.Errorf("创建临时文件失败：%w", err)
	}
	fileCloseHandled := false // 是否已在正常流程中调用并检查 Close。
	defer func() {
		if fileCloseHandled {
			return
		}
		// 提前返回时兜底关闭；正常流程已关闭则不再重复关闭。
		if err := outputFile.Close(); err != nil {
			fmt.Printf("关闭文件失败：%v\n", err)
		}
	}()
	// TODO 6：行号从 1 开始，写入“行号 + 空格 + 原始行”，检查写入错误。
	// 空文件不编号，空行要编号；末行无换行时先处理数据再判断 EOF。
	// 循环结束后检查原始字节数是否等于声明大小。
	body := &io.LimitedReader{R: reader, N: int64(size)}
	bodyReader := bufio.NewReader(body)
	lineNumber := 1
	bytesRead := 0
	for bytesRead < size {
		line, err := bodyReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取正文失败：%w", err)
		}
		bytesRead += len(line)
		if line != "" {
			_, writeErr := outputFile.WriteString(fmt.Sprintf("%d %s", lineNumber, line))
			if writeErr != nil {
				return fmt.Errorf("写入输出文件失败：%w", writeErr)
			}
			lineNumber++
		}
		if err == io.EOF {
			break
		}
	}

	// TODO 7：确认文件写入、关闭成功，取得该输出文件的新大小。
	// 出错返回 fmt.Errorf，不沿用 PA5 中会 panic 的 check。
	// 回复 "original: ... bytes, new: ... bytes\n"，检查写入和 Flush。
	if bytesRead != size {
		return fmt.Errorf("上传结束时字节数不匹配：读取 %d，声明 %d", bytesRead, size)
	}
	fmt.Printf("[%s] 客户端=%s 输出=%s\n",
		time.Now().Format("15:04:05.000"),
		conn.RemoteAddr(),
		outputFile.Name(),
	)

	// TODO 8：成功 Flush 后记录时间，再 Sleep(processingDelay)。
	// 延迟结束才返回 nil，由 defer 关闭 conn。不要把延迟放到回复之前。
	closeErr := outputFile.Close()
	fileCloseHandled = true
	if closeErr != nil {
		return fmt.Errorf("关闭输出文件失败：%w", closeErr)
	}

	info, err := os.Stat(outputFile.Name())
	if err != nil {
		return fmt.Errorf("获取输出大小失败：%w", err)
	}

	writer := bufio.NewWriter(conn)
	_, err = fmt.Fprintf(
		writer,
		"original: %d bytes, new: %d bytes\n",
		size,
		info.Size(),
	)
	if err != nil {
		return fmt.Errorf("写入回复失败：%w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("发送回复失败：%w", err)
	}

	time.Sleep(processingDelay)

	return nil
}
