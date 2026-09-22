package main

import (
	"bufio"
	"fmt"
	"net"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Launching server...")
	listener, _ := net.Listen("tcp", "127.0.0.1:8080")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		check(err)
		defer conn.Close()

		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		check(err)
		fmt.Printf("%s", message)

		writer := bufio.NewWriter(conn)
		newline := fmt.Sprintf("%d bytes received\n", len(message))
		_, err = writer.WriteString(newline)
		check(err)
		writer.Flush()
	}

}
