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
	//inputFlag = flag.String("", "stdin", "Input type")
	//byteFlag  = flag.String("c", "no_init", "count the number of bytes in a file")
	//lineFlag  = flag.String("l", "no_init", "count the number of lines in a file")
	//wordFlag  = flag.String("w", "no_init", "count the number of words in a file")
	inputFlag = flag.Bool("", false, "Input type")
	byteFlag  = flag.Bool("c", false, "count the number of bytes in a file")
	lineFlag  = flag.Bool("l", false, "count the number of lines in a file")
	wordFlag  = flag.Bool("w", false, "count the number of words in a file")
)

// PanicOnError panics if the error is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// FileExists checks if the file exists.
func FileExists(fileName string) bool {
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
func CountEntity(file *os.File, processor func([]byte) int) int {
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
func countBytes(file *os.File) int {
	return CountEntity(file, func(buffer []byte) int {
		return len(buffer)
	})
}

// countLines counts the number of lines in a file.
func countLines(file *os.File) int {
	return CountEntity(file, func(buffer []byte) int {
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
func countWords(file *os.File) int {
	processingWord := false
	return CountEntity(file, func(buffer []byte) int {
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

// GetTargetFile opens the required file for processing.
func GetTargetFile(fileName string) (file *os.File, err error) {
	// check if the input is from a file & the file exists.
	// higher priority is given to the file (than stdin)
	if len(fileName) > 0 {
		if !FileExists(fileName) {
			os.Exit(1)
		}
		return os.Open(fileName)
	}

	// check if the input is from stdin
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return os.Stdin, nil
	}
	return nil, fmt.Errorf("no input file specified")
}

func main() {
	/*
		Flags
			-c : count the number of bytes in a file
			-l : count the number of lines in a file
			-w : count the number of words in a file
	*/
	flag.Parse()

	file, err := GetTargetFile(flag.Arg(0))
	PanicOnError(err)
	defer func() {
		PanicOnError(file.Close())
	}()

	answer := -1
	if *byteFlag {
		answer = countBytes(file)
	} else if *lineFlag {
		answer = countLines(file)
	} else if *wordFlag {
		answer = countWords(file)
	} else {
		fmt.Println("No flag specified")
		os.Exit(1)
	}
	fmt.Println(answer)
}
