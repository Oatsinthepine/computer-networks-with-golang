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

// PA3.go currently connects to this local address by default.
var serverAddress = "127.0.0.1:12000"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 1. PA4 is a server: listen first, then accept one client.
	listener, err := net.Listen("tcp", serverAddress)
	check(err)
	defer listener.Close()
	fmt.Println("Waiting for an upload on", serverAddress)

	conn, err := listener.Accept()
	check(err)
	defer conn.Close()

	// 2. The first line from PA3 is the original file size, for example "64\n".
	reader := bufio.NewReader(conn)
	sizeLine, err := reader.ReadString('\n')
	check(err)
	originalSize, err := strconv.Atoi(strings.TrimSpace(sizeLine))
	check(err)
	if originalSize < 0 || originalSize > 1024*1024 {
		fmt.Println("Invalid file size: expected 0 to 1048576 bytes.")
		return
	}

	// 3. Create the transformed output file in the current working directory.
	outputFile, err := os.Create("./whatever.txt")
	check(err)

	// The LimitedReader makes sure the body reader never reads past this upload.
	// Keep one bufio.Reader for the whole body: it may already buffer later bytes.
	body := &io.LimitedReader{R: reader, N: int64(originalSize)}
	print(body)
	bodyReader := bufio.NewReader(body)
	bytesRead := 0
	lineNumber := 1

	// 4. Read one text line at a time, number it, then write it to whatever.txt.
	for bytesRead < originalSize {
		line, err := bodyReader.ReadString('\n')
		if err != nil && err != io.EOF {
			check(err)
		}
		bytesRead += len(line) // len(string) counts bytes, matching PA3's file size.
		_, writeErr := fmt.Fprintf(outputFile, "%d %s", lineNumber, line)
		check(writeErr)
		lineNumber++

		// EOF is normal for a final line without a trailing newline.
		if err == io.EOF {
			break
		}
	}

	if bytesRead != originalSize {
		check(fmt.Errorf("upload ended after %d of %d bytes", bytesRead, originalSize))
	}

	// 5. The new size includes the line numbers and spaces we added.
	check(outputFile.Close())
	info, err := os.Stat("./whatever.txt")
	check(err)

	// 6. Send one application-level reply line back to PA3.
	writer := bufio.NewWriter(conn)
	_, err = fmt.Fprintf(writer, "original: %d bytes, new: %d bytes\n", originalSize, info.Size())
	check(err)
	check(writer.Flush())
	// 7. main returns, so the deferred connection and listener closes run.
}
