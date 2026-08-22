package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
冗余代码：check_file_exists 没必要
在 Go 中，不需要先检查文件是否存在再去 os.Open。
因为 os.Open(filename) 本身就会返回 err。先 os.Stat 再 os.Open 不仅多写了代码，还引入了竞态条件（比如检查时文件在，但在 Open 之前的瞬间被删了）。

直接 os.Open，然后判断 err 即可。
*/

//func check_file_exists(filename string) bool {
//	_, err := os.Stat(filename)
//	return !os.IsNotExist(err)
//}

func main() {
	// 1. prompt the user for the input and output file names
	fmt.Printf("Please enter the input and output file names, separate using space: ")
	var inputFilename, outputFilename string
	fmt.Scanf("%s %s", &inputFilename, &outputFilename)

	// 2. reads from the input file one line at time
	// if !check_file_exists(inputFilename) {
	// 	fmt.Printf("Input file %s does not exist. Exiting program.\n", inputFilename)
	// 	return
	// }

	inputFile, err := os.Open(inputFilename)

	if err != nil {
		fmt.Printf("Error opening input file %s: %v. Exiting program.\n", inputFilename, err)
		return
	}

	defer inputFile.Close()

	// 3. prepends the line count to each line and writes the line into the output file
	scanner := bufio.NewScanner(inputFile)

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		line_count := fmt.Sprintf("%d", len(line))
		_, err := writer.WriteString(fmt.Sprintf("%d %s %s\n", lineNumber, line_count, line))

		if err != nil {
			panic(err)
		}

		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println("File processing completed successfully.")

}
