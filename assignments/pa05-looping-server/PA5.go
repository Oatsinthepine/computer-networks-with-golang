package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const serverAddress = "127.0.0.1:12000"
const maxFileSize = 1 << 20 // 1 MiB，按字节限制原文件大小。

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Launched PA5 服务端。")

	// TODO 1：在循环外调用 net.Listen，检查错误。
	// 监听失败时打印原因并退出；成功后安排 listener.Close()。
	listener, err := net.Listen("tcp", serverAddress)
	check(err)
	defer listener.Close()

	// TODO 2：使用 for 循环反复调用 listener.Accept()，检查错误。
	// Accept 失败时打印原因并退出，避免在失效的 listener 上空转。
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受客户端连接失败：", err)
			return
		}

		// TODO 3：同步调用 handleUpload(conn)，不要加 go。
		// 本次上传失败时打印错误，然后继续接收下一位客户端。
		// 不要在这个无限循环里 defer conn.Close()。
		if err := handleUpload(conn); err != nil {
			fmt.Println("处理上传失败：", err)
		}
	}
}

func handleUpload(conn net.Conn) error {
	// TODO 4：在这个函数中安排 defer conn.Close()。
	// 每次函数返回，就结束当前连接；listener 仍然存在。
	defer conn.Close()
	// TODO 5：建立 bufio.Reader，读大小行，去掉换行后转为整数。
	// 检查大小在 0 到 maxFileSize 之间；失败时返回 error。
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
	// TODO 6：创建 whatever.txt（覆盖旧输出），安排文件关闭。
	// 输出路径相对于程序的运行目录；不要覆盖 testdata/input.txt。
	outputFile, err := os.Create("whatever.txt")
	if err != nil {
		return fmt.Errorf("创建输出文件失败：%w", err)
	}
	defer outputFile.Close()
	// TODO 7：用原来的 reader 建立有限正文读取范围，再按行读取。
	// 每次上传都重新初始化行号为 1、原始字节计数为 0。
	// 保留原始换行，在非空的读取结果前添加“行号 + 空格”。
	// 空行 "\n" 也要编号；空文件不应生成行号。
	// 最后一行没有换行时，先处理返回的数据，再判断 EOF。
	// 若实际读取字节数少于声明大小，返回错误，不回复成功。
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

	if bytesRead != size {
		return fmt.Errorf("上传结束时字节数不匹配：读取 %d，声明 %d", bytesRead, size)
	}

	// TODO 8：确认输出写入和关闭成功，取得新文件字节大小。
	// 回复一行 "original: ... bytes, new: ... bytes\n"，检查 Flush。
	// 成功完成后才返回 nil；过程中出错返回 error，不使用 panic。
	check(outputFile.Close())
	info, err := os.Stat("./whatever.txt")
	check(err)

	writer := bufio.NewWriter(conn)
	_, err = fmt.Fprintf(writer, "original: %d bytes, new: %d bytes\n", size, info.Size())
	check(err)
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("发送回复失败：%w", err)
	}

	return nil
}
