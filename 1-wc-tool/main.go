package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Flags
var (
	byteFlag = flag.Bool("c", false, "count bytes")
	lineFlag = flag.Bool("l", false, "count lines")
	wordFlag = flag.Bool("w", false, "count words")
	charFlag = flag.Bool("m", false, "count characters")
)

// PanicOnError panics if the error is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type Output struct {
	lines int
	words int
	bytes int
	chars int
}

func (op *Output) String() string {
	var outStr []string
	if *lineFlag {
		outStr = append(outStr, fmt.Sprintf("%d", op.lines))
	}
	if *wordFlag {
		outStr = append(outStr, fmt.Sprintf("%d", op.words))
	}
	if *charFlag {
		outStr = append(outStr, fmt.Sprintf("%d", op.chars))
	}
	if *byteFlag {
		outStr = append(outStr, fmt.Sprintf("%d", op.bytes))
	}
	return strings.Join(outStr, "\t")
}

// CountEntity counts the number of lines, words, bytes and characters from the reader.
// It returns the Output struct.
func CountEntity(reader *bufio.Reader) Output {
	op := Output{}
	var prevChar rune
	for {
		chRead, bytesRead, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if prevChar != rune(0) && !unicode.IsSpace(prevChar) {
					op.words++
				}
				break
			}
			PanicOnError(err)
		}
		op.bytes += bytesRead
		op.chars++
		if chRead == '\n' {
			op.lines++
		}
		// if prev is not space and current is space, then it is a word
		if !unicode.IsSpace(prevChar) && unicode.IsSpace(chRead) {
			op.words++
		}
		prevChar = chRead
	}
	return op
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

func main() {
	flag.Parse()

	file, err := GetTargetFile(flag.Arg(0))
	PanicOnError(err)
	defer func() {
		PanicOnError(file.Close())
	}()

	// If no option is set, use all options.
	if !(*byteFlag || *lineFlag || *wordFlag || *charFlag) {
		*byteFlag = true
		*lineFlag = true
		*wordFlag = true
		*charFlag = true
	}

	reader := bufio.NewReader(file)
	op := CountEntity(reader)
	fmt.Println(op.String())
}
