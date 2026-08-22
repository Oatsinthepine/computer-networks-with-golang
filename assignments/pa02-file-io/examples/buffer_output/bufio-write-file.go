package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("PA2-output.txt")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	len, _ := writer.WriteString("This is a test string.\n")
	fmt.Println(len)

	writer.Flush()
}
