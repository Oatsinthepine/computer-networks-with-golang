package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Open("hello-world.go")
	check(err)

	defer f.Close()

	word1, word2 := "", ""
	fmt.Fscanln(f, &word1, &word2)
	fmt.Printf("Word 1: %s\nWord 2: %s\n", word1, word2)

	for i := 1; i < 6; i++ {
		word1, word2 = "", ""
		fmt.Fscanln(f, &word1, &word2)
		fmt.Printf("Word 1: %s\nWord 2: %s\n", word1, word2)
	}
}
