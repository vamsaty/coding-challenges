package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Flags
var (
	inputFlag = flag.String("", "stdin", "Input type")
	byteFlag  = flag.String("c", "no_init", "count the number of bytes in a file")
	lineFlag  = flag.String("l", "no_init", "count the number of lines in a file")
	wordFlag  = flag.String("w", "no_init", "count the number of words in a file")
)

// PanicOnError panics if the error is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// IsValidFlag checks if the file exists.
func IsValidFlag(fileName string) bool {
	if fileName == "no_init" {
		return false
	}
	_, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("File %s does not exist\n", fileName)
		return false
	}
	return true
}

// CountEntity is a common function for counting bytes, lines, and words.
// It makes use of a `processor` function to process the buffer and count
// the occurrence of the required entity. The steps were common for
// counting bytes, lines, and words.
func CountEntity(fileName string, processor func([]byte) int) int {
	// open the file for reading
	file, err := os.Open(fileName)
	PanicOnError(err)
	defer func() {
		PanicOnError(file.Close())
	}()
	// read the file in as a stream of bytes and process the read data
	counter := 0
	for {
		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// process the buffer
				counter += processor(buffer[:n])
				return counter
			}
			PanicOnError(err)
		}
		// process the buffer
		counter += processor(buffer[:n])
	}
}

// countBytes counts the number of bytes in a file.
func countBytes(fileName string) int {
	return CountEntity(fileName, func(buffer []byte) int {
		return len(buffer)
	})
}

// countLines counts the number of lines in a file.
func countLines(fileName string) int {
	return CountEntity(fileName, func(buffer []byte) int {
		count := 0
		for i := 0; i < len(buffer); i++ {
			if rune(buffer[i]) == '\n' {
				count++
			}
		}
		return count
	})
}

// countWords counts the number of words in a file.
func countWords(fileName string) int {
	processingWord := false
	return CountEntity(fileName, func(buffer []byte) int {
		count := 0
		for i := 0; i < len(buffer); i++ {
			if unicode.IsSpace(rune(buffer[i])) {
				if processingWord {
					count++
				}
				processingWord = false
			} else {
				processingWord = true
			}
		}
		return count
	})
}

func main() {
	/*
		Flags
			-c : count the number of bytes in a file
			-l : count the number of lines in a file
			-w : count the number of words in a file
	*/
	flag.Parse()
	answer := -1

	if IsValidFlag(*byteFlag) {
		answer = countBytes(*byteFlag)
	} else if IsValidFlag(*lineFlag) {
		answer = countLines(*lineFlag)
	} else if IsValidFlag(*wordFlag) {
		answer = countWords(*wordFlag)
	} else {
		fmt.Println("Invalid flag")
		os.Exit(1)
	}
	fmt.Println(answer)
}
