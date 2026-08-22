package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	f, err := os.Open("hello-world.go")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}
