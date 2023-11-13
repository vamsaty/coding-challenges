package server

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
)

/* ---------------- HttpServerConfig ---------------- */

// HttpServerConfig is a struct that contains the host and port of the server.
type HttpServerConfig struct {
	Host string
	Port int
}

func (sc *HttpServerConfig) String() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

/* ---------------- HttpServer ---------------- */

// HttpServer is an interface that defines the methods of an HTTP server.
type HttpServer interface {
	StartServer()
	Init()
	Register(method, path string, Handler MyHttpHandler)
	StopServer()
}

func NewHttpServer(config *HttpServerConfig) HttpServer {
	logger, _ := zap.NewProduction()
	return &SimpleServerImpl{
		Logger:           logger,
		HttpServerConfig: config,
		Once:             &sync.Once{},
		Registry:         NewRegistry(),
	}
}

/* ---------------- SimpleServerImpl ---------------- */

// SimpleServerImpl satisfies the HttpServer interface.
type SimpleServerImpl struct {
	*sync.Once
	*zap.Logger
	*HttpServerConfig
	*Registry
	Listener *net.TCPListener
}

// Register registers a handler for a given method and path.
func (server *SimpleServerImpl) Register(method, path string, Handler MyHttpHandler) {
	server.Registry.Register(method, path, Handler)
}

// StartServer starts the HTTP server after initialization
func (server *SimpleServerImpl) StartServer() {
	server.Init()
	server.Logger.Info("server started")
	for {
		conn, err := server.Listener.AcceptTCP()
		if err != nil {
			server.Logger.Error("failed to accept connection", zap.Error(err))
			continue
		}
		go server.handleConnection(conn)
	}
}

// Init initializes the HTTP server.
func (server *SimpleServerImpl) Init() {
	server.Logger.Info("starting server")
	server.Once.Do(func() {
		listener, err := net.Listen("tcp", server.HttpServerConfig.String())
		if err != nil {
			server.Logger.Fatal("failed to start server", zap.Error(err))
		}
		server.Listener = listener.(*net.TCPListener)
	})
}

// handleConnection handles a single connection
func (server *SimpleServerImpl) handleConnection(conn *net.TCPConn) {
	connId := zap.String("remote", conn.RemoteAddr().String())
	server.Logger.Info("handling connection.", connId)

	var err error
	defer func() {
		server.Logger.Info("closing connection.", connId)
		if err = conn.Close(); err != nil {
			server.Logger.Error("failed to close connection", zap.Error(err))
		}
	}()

	reader := NewRequestParser(conn)
	writer := NewResponseWriter(conn)
	for {
		server.Logger.Info("parsing request", connId)
		if err = reader.Parse(); err != nil {
			if err == io.EOF {
				server.Logger.Info("connection closed by client", connId)
				return
			}
			server.Logger.Error("failed to parse request", zap.Error(err))
			continue
		}
		req := reader.GetRequest()
		server.Logger.Info("request parsed", connId, zap.String("method", req.Method), zap.String("path", req.Path))
		RequestHandler := server.Registry.GetHandler(req.Method, req.Path)
		RequestHandler(writer, reader.Request)
	}
}

// StopServer stops the HTTP server.
func (server *SimpleServerImpl) StopServer() {
	if server.Listener != nil {
		err := server.Listener.Close()
		server.Logger.Fatal("failed to stop server", zap.Error(err))
	}
	server.Logger.Info("server stopped")
}
