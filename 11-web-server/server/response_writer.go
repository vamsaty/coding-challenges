package server

import (
	"fmt"
	"net"
	"time"
)

/* ---------------- ResponseWriter ---------------- */

// ResponseWriter is an interface to write the response to the client.
type ResponseWriter interface {
	Write(data *Response) error
}

/* ---------------- SimpleResponseWriter ---------------- */

// SimpleResponseWriter is vanilla implementation of ResponseWriter.
type SimpleResponseWriter struct {
	net.Conn
}

func NewResponseWriter(conn net.Conn) *SimpleResponseWriter {
	return &SimpleResponseWriter{Conn: conn}
}

// Write is just a wrapper around net.Conn.Write.
func (srw *SimpleResponseWriter) Write(data *Response) (err error) {
	fmt.Println("SimpleResponseWriter.Write", data)
	err = srw.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	fmt.Println("ERR", err)
	_, err = srw.Conn.Write(data.Serialize())
	return err
}
