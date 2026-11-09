package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

const documentRoot = "testdata"

const serverAddress = "127.0.0.1:12000"

func main() {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("PA7 listening on", serverAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 先按顺序处理，每个请求由处理函数负责关闭连接。
		if err := serveRequest(conn, documentRoot, os.Stdout); err != nil {
			fmt.Println("Request error:", err)
		}
	}
}

// serveRequest 每条连接只处理一个 GET 请求，不发送 HTTP 响应。
func serveRequest(conn net.Conn, rootDir string, output io.Writer) error {
	defer conn.Close()
	// 避免未发送完整 headers 的客户端永远占用顺序服务端。
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}
	reader := bufio.NewReader(io.LimitReader(conn, 16*1024))

	line, err := reader.ReadString('\n')

	if err != nil {
		return fmt.Errorf("读取请求行：%w", err)
	}
	if !strings.HasSuffix(line, "\r\n") {
		return fmt.Errorf("请求行需要 CRLF")
	}
	parts := strings.Split(strings.TrimSuffix(line, "\r\n"), " ")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("请求行需要三个字段")
	}
	if parts[2] != "HTTP/1.1" && parts[2] != "HTTP/1.0" {
		return fmt.Errorf("仅支持 HTTP/1.0 和 HTTP/1.1")
	}

	// headers 以空行结束，不能等连接 EOF。
	host := ""
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取 headers：%w", err)
		}
		if line == "\r\n" {
			break
		}
		if !strings.HasSuffix(line, "\r\n") {
			return fmt.Errorf("header 需要 CRLF")
		}
		name, value, ok := strings.Cut(strings.TrimSuffix(line, "\r\n"), ":")
		if !ok || name == "" || strings.ContainsAny(name, " \t") {
			return fmt.Errorf("无效 header")
		}
		if strings.EqualFold(name, "Host") {
			if host != "" {
				return fmt.Errorf("重复 Host")
			}
			host = strings.TrimSpace(value)
		}
	}
	if parts[2] == "HTTP/1.1" && host == "" {
		return fmt.Errorf("HTTP/1.1 缺少 Host")
	}
	if parts[0] != "GET" {
		return fmt.Errorf("本练习只支持 GET")
	}
	fmt.Fprintf(output, "Method: %s, Target: %s, Host: %s\n", parts[0], parts[1], host)
	return printFile(rootDir, parts[1], output)
}

// printFile 将 URL 路径限制在指定目录内，打印文件信息而不读取正文。
func printFile(rootDir, target string, output io.Writer) error {
	u, err := url.ParseRequestURI(target)
	if err != nil || !strings.HasPrefix(target, "/") || strings.HasPrefix(target, "//") {
		return fmt.Errorf("无效的请求目标：%q", target)
	}
	path := u.Path // 已解码的路径，不包含查询参数。
	if path == "/" {
		path = "/hello.html"
	}
	for _, part := range strings.Split(path, "/") {
		if part == ".." {
			return fmt.Errorf("不允许访问上级目录")
		}
	}
	if strings.ContainsAny(path, "\\\x00") {
		return fmt.Errorf("无效的文件路径")
	}

	// Go 的 Root 将文件访问限制在目录内，也限制符号链接逃出目录。
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return fmt.Errorf("打开网站目录：%w", err)
	}
	defer root.Close()
	info, err := root.Stat(strings.TrimPrefix(path, "/"))
	if os.IsNotExist(err) {
		fmt.Fprintln(output, "File not found")
		return nil
	}
	if err != nil {
		return fmt.Errorf("查询文件：%w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("请求目标不是普通文件")
	}
	fmt.Fprintf(output, "File: %s, Size: %d bytes\n", path, info.Size())
	return nil
}
