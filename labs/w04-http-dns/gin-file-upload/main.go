package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize  = int64(1 << 20) // 1 MiB file limit
	maxRequestSize = maxUploadSize + 64<<10
)

type uploadServer struct {
	outputDir string
}

func main() {
	router := newRouter("outputs")
	fmt.Println("Gin upload server: http://127.0.0.1:8080")
	check(router.Run("127.0.0.1:8080"))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func newRouter(outputDir string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	// This lab only serves localhost and does not sit behind a reverse proxy.
	check(router.SetTrustedProxies(nil))

	server := uploadServer{outputDir: outputDir}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/upload", server.upload)

	return router
}

func (s uploadServer) upload(c *gin.Context) {
	// Limit the complete HTTP request before parsing multipart/form-data.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSize)
	// 寻找 multipart表单中名字为 file 的部分
	fileHeader, err := c.FormFile("file")

	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "upload request is too large"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "multipart form field 'file' is required"})
		return
	}
	if fileHeader.Size > maxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file exceeds the 1 MiB limit"})
		return
	}

	source, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open uploaded file"})
		return
	}
	defer source.Close()

	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot prepare output directory"})
		return
	}

	// Generate the name on the server instead of trusting the client filename.
	output, err := os.CreateTemp(s.outputDir, "whatever-*.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create output file"})
		return
	}
	outputPath := output.Name()

	originalSize, newSize, transformErr := addLineNumbers(source, output)
	closeErr := output.Close()
	if transformErr != nil || closeErr != nil {
		_ = os.Remove(outputPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot transform uploaded file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "upload transformed successfully",
		"original_size": originalSize,
		"new_size":      newSize,
		"output":        filepath.Base(outputPath),
	})
}

// addLineNumbers contains the same business logic as PA4, independent of Gin.
func addLineNumbers(source io.Reader, destination io.Writer) (int64, int64, error) {
	reader := bufio.NewReader(source)
	var originalSize int64
	var newSize int64

	for lineNumber := 1; ; lineNumber++ {
		line, readErr := reader.ReadString('\n')
		if len(line) > 0 {
			originalSize += int64(len(line))
			written, writeErr := fmt.Fprintf(destination, "%d %s", lineNumber, line)
			if writeErr != nil {
				return originalSize, newSize, writeErr
			}
			newSize += int64(written)
		}

		if readErr == io.EOF {
			return originalSize, newSize, nil
		}
		if readErr != nil {
			return originalSize, newSize, readErr
		}
	}
}
