package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const documentRoot = "../native/testdata"

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

func serveRequest(conn net.Conn, rootDir string, output io.Writer) error {
	defer conn.Close()
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}
	reader := bufio.NewReader(io.LimitReader(conn, 16*1024))

	// 标准库代替手动解析请求行和 headers；这里仍自行管理 TCP 连接。
	req, err := http.ReadRequest(reader)
	if err != nil {
		return fmt.Errorf("解析 HTTP 请求：%w", err)
	}

	// 本练习不读取请求正文，完成后直接关闭连接。
	if req.Method != "GET" {
		return fmt.Errorf("本练习只支持 GET")
	}
	if req.Proto != "HTTP/1.0" && req.Proto != "HTTP/1.1" {
		return fmt.Errorf("不支持的 HTTP 版本")
	}
	if req.Proto == "HTTP/1.1" && req.Host == "" {
		return fmt.Errorf("缺少 Host")
	}
	fmt.Fprintf(output, "Method: %s, Target: %s, Host: %s\n", req.Method, req.RequestURI, req.Host)
	return printFile(rootDir, req.RequestURI, output)
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
