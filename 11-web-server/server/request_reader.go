package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

var (
	supportedMethods = []string{"GET", "POST", "PATCH", "DELETE", "PUT", "CONNECT", "OPTIONS", "TRACE"}
	httpVersionRegex = regexp.MustCompile(`HTTP\/\d\.\d`)
	crlfRegex        = regexp.MustCompile(`^\r\n$`)
)

/* ---------------- RequestParser ---------------- */

type RequestParser interface {
	GetRequest() Request
	Parse() error
}

func NewRequestParser(conn net.Conn) *SimpleRequestParser {
	return &SimpleRequestParser{
		Reader:  bufio.NewReader(conn),
		Request: &Request{},
	}
}

/* ---------------- SimpleRequestParser ---------------- */

// SimpleRequestParser is a simple implementation of RequestParser.
// It parses the byte slice into a Request. (see Parse method)
type SimpleRequestParser struct {
	*bufio.Reader
	*Request
}

// GetRequest returns a copy of the Request.
func (srp *SimpleRequestParser) GetRequest() Request { return srp.Request.Clone() }

// Parse parses the request line, header, and body.
func (srp *SimpleRequestParser) Parse() error {
	var err error
	if err = srp.ParseRequestLine(); err != nil {
		return err
	}
	if err = srp.ParseHeader(); err != nil {
		return err
	}
	if err = srp.ParseBody(); err != nil {
		return err
	}
	return nil
}

// ParseRequestLine parses the request line.
func (srp *SimpleRequestParser) ParseRequestLine() error {
	if line, _, err := srp.ReadLine(); err != nil {
		return err
	} else {
		parts := strings.Split(string(line), " ")
		if len(parts) != 3 {
			return fmt.Errorf("invalid request line: %s", string(line))
		}
		srp.Request.RequestLine = RequestLine{
			Method:  parts[0],
			Path:    parts[1],
			Version: parts[2],
		}
	}
	return srp.RequestLine.validate()
}

// ParseHeader parses the request header and sets the Header field in Request.
// Reads a line and parses the key-value pair, until a blank line is found
// (line starting with CRLF)
func (srp *SimpleRequestParser) ParseHeader() error {
	header := Header{data: make(map[string]string)}
	for {
		if line, _, err := srp.ReadLine(); err != nil {
			return err
		} else if crlfRegex.Match(line) || len(line) == 0 {
			break
		} else {
			parts := strings.Split(string(line), ":")
			header.Add(parts[0], parts[1])
		}
	}
	srp.Request.Header = header
	return nil
}

// ParseBody parses the body. Reads until the "Content-Length" is reached.
// In case of GET methods (the body is empty), it returns nil
// (ContentLength is assumed to be 0).
func (srp *SimpleRequestParser) ParseBody() error {
	body := Body{Data: []byte{}}

	// Keep appending the payload into the Body
	for body.Len() < srp.Request.ContentLength() {
		if line, _, err := srp.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else {
			// append line to body.Data
			body.Data = append(body.Data, line...)
		}
	}
	srp.Request.Body = body
	return nil
}
