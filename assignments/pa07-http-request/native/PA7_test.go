package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRequests(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "hello.html"), []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "nested"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "nested", "page.html"), []byte("nested"), 0600); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "secret.html")
	if err := os.WriteFile(outside, []byte("private"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape.html")); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, target, want string
		bad                bool
	}{
		{"exists", "/hello.html", "Size: 5 bytes", false},
		{"root", "/", "Size: 5 bytes", false},
		{"query", "/hello.html?x=1", "Size: 5 bytes", false},
		{"nested", "/nested/page.html", "Size: 6 bytes", false},
		{"missing", "/missing.html", "File not found", false},
		{"directory", "/nested/", "", true},
		{"traversal", "/%2e%2e/secret.html", "", true},
		{"symlink", "/escape.html", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			client, server := net.Pipe()
			client.SetDeadline(time.Now().Add(2 * time.Second))
			defer client.Close()
			done := make(chan error, 1)
			go func() { done <- serveRequest(server, root, &out) }()
			_, err := io.WriteString(client, "GET "+tc.target+" HTTP/1.1\r\nHost: localhost\r\n\r\n")
			if err != nil {
				t.Fatal(err)
			}
			response, err := io.ReadAll(client)
			if err != nil {
				t.Fatal(err)
			}
			if len(response) != 0 {
				t.Fatalf("PA7 should not send a response: %q", response)
			}
			err = <-done
			if (err != nil) != tc.bad {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.bad && !strings.Contains(out.String(), tc.want) {
				t.Fatalf("log: %s", &out)
			}
		})
	}
}
func TestRejectedRequests(t *testing.T) {
	for _, request := range []string{
		"BROKEN\r\n\r\n",
		"POST /hello.html HTTP/1.1\r\nHost: localhost\r\n\r\n",
		"GET /hello.html HTTP/1.1\r\nBrokenHeader\r\n\r\n",
	} {
		client, server := net.Pipe()
		client.SetDeadline(time.Now().Add(2 * time.Second))
		done := make(chan error, 1)
		go func() { done <- serveRequest(server, t.TempDir(), io.Discard) }()
		// 解析器可能在读取整条请求之前拒绝并关闭连接。
		io.WriteString(client, request)
		client.Close()
		if err := <-done; err == nil {
			t.Fatalf("accepted %q", request)
		}
	}
}
