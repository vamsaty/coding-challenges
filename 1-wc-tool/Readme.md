# Write Your Own wc Tool

## Description
The challenge was to use build your own wc tool. The tool should be able to take in a file or stdin and output the number of lines, words, and characters in the file or stdin.

## Usage

Steps to build the binary and execute it -
```
go run ./main.go -c <filename>
```
or 
```
cat <filename> | go run ./main.go -c
```

## Flags

| Flag | Description | Default |
| --- | --- | --- |
| -c | Count the number of characters | false |
| -l | Count the number of lines | false |
| -w | Count the number of words | false |

#### NOTE: when both a file and stdin are provided, the file will be used.