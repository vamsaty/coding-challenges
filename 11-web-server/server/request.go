package server

import "fmt"

/* ---------------- Request ---------------- */

// Request represents an HTTP request, consisting of RequestLine, Header and Body.
type Request struct {
	RequestLine // start line request
	Header      // header
	Body        // body
}

// Clone returns a copy of the Request.
func (req *Request) Clone() Request {
	return Request{
		RequestLine: req.RequestLine.Clone(),
		Header:      req.Header.Clone(),
		Body:        req.Body.Clone(),
	}
}

/* ---------------- RequestLine ---------------- */

// RequestLine is the start line of an HTTP request
type RequestLine struct {
	Method  string
	Path    string
	Version string
}

// Clone returns a copy of the RequestLine.
func (rl *RequestLine) Clone() RequestLine {
	return RequestLine{
		Method:  rl.Method,
		Path:    rl.Path,
		Version: rl.Version,
	}
}

// validate checks if the fields in RequestLine are valid
func (rl *RequestLine) validate() error {
	found := false
	for _, m := range supportedMethods {
		if m == rl.Method {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("unsupported method %s", rl.Method)
	}
	if !httpVersionRegex.Match([]byte(rl.Version)) {
		return fmt.Errorf("unsupported http version %s", rl.Version)
	}
	return nil
}
