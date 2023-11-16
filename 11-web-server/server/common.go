package server

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// MyHttpHandler is a function that handles HTTP requests.
// An HTTP server registers various MyHttpHandler functions to handle
// different requests, depending on request METHOD and URI.
type MyHttpHandler func(writer ResponseWriter, request *Request)

/* ---------------- START : Common Handlers ---------------- */
var (
	notFoundResponse = func() *Response {
		return NewResponseBuilder().StatusCode(404).Phrase("Not Found").
			H("Content-Type", "text/html").Body([]byte("the file doesn't exist")).Build()
	}
	internalServerErrorResponse = func() *Response {
		return NewResponseBuilder().StatusCode(500).Phrase("Internal Server Error").
			H("Content-Type", "text/html").Body([]byte("internal server error")).Build()
	}
)

func Push(writer ResponseWriter, builder *ResponseBuilder) {
	PushResponse(writer, builder.Build())
}

func PushResponse(writer ResponseWriter, response *Response) {
	if err := writer.Write(response); err != nil {
		fmt.Println("failed to write response", err)
	}
}

// HandlerNotFound is the default handler if no handler exists
// for the requested URL and Method.
func HandlerNotFound(writer ResponseWriter, _ *Request) {
	PushResponse(writer, notFoundResponse())
}

// FetchWebPage reads the file from the local file system and sends it in the response body.
// Files only in the "./www" directory are accessible.
func FetchWebPage(writer ResponseWriter, request *Request) {
	rb := NewResponseBuilder()
	location := "./www" + request.Path
	fmt.Println("location", location)
	file, err := os.OpenFile(location, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("failed to open file", err)
		PushResponse(writer, notFoundResponse())
		return
	}
	defer file.Close()

	// send the file data in the response body
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("failed to read file", err)
		PushResponse(writer, internalServerErrorResponse())
		return
	}

	rb = rb.StatusCode(200).
		H("Content-Type", "text/html").
		H("Content-Length", strconv.Itoa(len(data))).
		Body(data)
	Push(writer, rb)
}

/* ---------------- END : Common Handlers ---------------- */

/* ---------------- START : Header ---------------- */

// Header represents the HTTP header section
type Header struct {
	data map[string]string
}

func NewHeader() Header {
	return Header{data: map[string]string{}}
}

// ContentLength returns the content length of the header
func (h *Header) ContentLength() int {
	if value, found := h.data["Content-Length"]; found {
		v, _ := strconv.ParseInt(value, 10, 32)
		return int(v)
	}
	return 0
}

// Get returns the value of the key
func (h *Header) Get(key string) string { return h.data[key] }

// Clone returns a copy of the header
func (h *Header) Clone() Header {
	header := Header{data: map[string]string{}}
	for key, value := range h.data {
		header.data[key] = value
	}
	return header
}

// Add adds a key-value pair to the header
func (h *Header) Add(key, value string) { h.data[key] = value }

// Serialize serializes the header into a byte slice
// following the HTTP protocol for HTTP header
func (h *Header) Serialize() []byte {
	sb := strings.Builder{}
	for key, value := range h.data {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	return []byte(sb.String())
}

/* ---------------- END : Header ---------------- */

/* ---------------- START : Body ---------------- */

// Body represents the HTTP body section
type Body struct {
	Data []byte
}

func NewBody() Body { return Body{Data: []byte{}} }

// Len returns the length of the body
func (body *Body) Len() int { return len(body.Data) }

// Clone returns a copy of the body
func (body *Body) Clone() Body {
	data := make([]byte, len(body.Data))
	copy(data, body.Data)
	return Body{Data: data}
}

// Serialize serializes the body into a byte slice
func (body *Body) Serialize() []byte { return body.Data }

/* ---------------- END : Body ---------------- */
