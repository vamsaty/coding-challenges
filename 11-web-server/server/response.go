package server

import "fmt"

/* ---------------- Response ---------------- */

// Response is the server's interpretation of the HTTP request.
// ResponseBuilder is used to build the response
type Response struct {
	statusLine StatusLine
	header     Header
	body       Body
}

// Clone returns a clone of the response
func (resp *Response) Clone() *Response {
	return &Response{
		statusLine: resp.statusLine.Clone(),
		header:     resp.header.Clone(),
		body:       resp.body.Clone(),
	}
}

// Serialize serializes the response into a byte array
func (resp *Response) Serialize() []byte {
	var data []byte
	data = append(data, resp.statusLine.Serialize()...)
	data = append(data, resp.header.Serialize()...)
	data = append(data, "\r\n"...)
	data = append(data, resp.body.Serialize()...)
	return data
}

/* ---------------- Status Line ---------------- */

// StatusLine represents the HTTP status line of an HTTP response.
// This is the start line of a response.
type StatusLine struct {
	Version    string
	StatusCode int
	Phrase     string
}

func NewDefaultStatusLine() StatusLine {
	return StatusLine{
		Version: "HTTP/1.1",
	}
}

// Clone returns a clone of the status line
func (sl *StatusLine) Clone() StatusLine {
	return StatusLine{
		Version:    sl.Version,
		StatusCode: sl.StatusCode,
		Phrase:     sl.Phrase,
	}
}

// Serialize serializes the status line into a byte array
func (sl *StatusLine) Serialize() []byte {
	return []byte(fmt.Sprintf("%s %d %s\r\n", sl.Version, sl.StatusCode, sl.Phrase))
}

/* ---------------- ResponseBuilder ---------------- */

// ResponseBuilder is a builder for HTTP Response
type ResponseBuilder struct {
	*Response
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{Response: &Response{
		statusLine: StatusLine{
			Version: "HTTP/1.1",
		},
		header: NewHeader(),
		body:   NewBody(),
	}}
}

// StatusCode sets the status code of the response
func (rb *ResponseBuilder) StatusCode(statusCode int) *ResponseBuilder {
	rb.statusLine.StatusCode = statusCode
	return rb
}

// Phrase sets the phrase of the response
func (rb *ResponseBuilder) Phrase(phrase string) *ResponseBuilder {
	rb.statusLine.Phrase = phrase
	return rb
}

// Version sets the HTTP version of the response
func (rb *ResponseBuilder) Version(version string) *ResponseBuilder {
	rb.statusLine.Version = version
	return rb
}

// H sets a header key-value pair
func (rb *ResponseBuilder) H(key, value string) *ResponseBuilder {
	rb.header.Add(key, value)
	return rb
}

// Data sets the body of the response
func (rb *ResponseBuilder) Data(data string) *ResponseBuilder {
	rb.body = Body{Data: []byte(data)}
	return rb
}

func (rb *ResponseBuilder) Body(data []byte) *ResponseBuilder {
	bodyData := make([]byte, len(data))
	copy(bodyData, data)
	rb.body = Body{Data: bodyData}
	return rb
}

// Build builds the response
func (rb *ResponseBuilder) Build() *Response {
	return rb.Response.Clone()
}
