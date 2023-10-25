package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, lines bool, countBytes bool) int {
	scanner := bufio.NewScanner(r)

	if !lines {
		scanner.Split(bufio.ScanWords)
	}

	if countBytes {
		bytesCounter := 0
		for scanner.Scan() {
			bytesCounter += len(scanner.Bytes())
		}
		return bytesCounter
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}

func main() {
	lines := flag.Bool("l", false, "Count Lines")
	byteCount := flag.Bool("b", false, "Count Bytes")

	flag.Parse()

	fmt.Println(count(os.Stdin, *lines, *byteCount))
}
