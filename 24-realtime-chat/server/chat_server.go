package server

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"net"
	"sync"
)

type ChatServer struct {
	clients map[string]net.Conn

	lock *sync.Mutex
	*zap.Logger
}

func NewChatServer() *ChatServer {
	logger, _ := zap.NewProduction()
	return &ChatServer{
		clients: make(map[string]net.Conn),
		Logger:  logger,
		lock:    &sync.Mutex{},
	}
}

func (s *ChatServer) close(conn net.Conn) error {
	return conn.Close()
}

func (s *ChatServer) addClient(conn net.Conn) (string, error) {
	// receive name
	s.lock.Lock()
	defer s.lock.Unlock()

	line, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		s.Warn(
			"error reading from connection.",
			zap.Error(err),
		)
		return "", err
	}
	newClient := string(line)
	s.Info("Added client.", zap.String("name", newClient))
	s.clients[newClient] = conn
	return newClient, nil
}

func (s *ChatServer) handleConnection(conn net.Conn) {
	var err error
	defer func() {
		s.Info("closing connection...")
		if err = s.close(conn); err != nil {
			s.Warn("error closing connection.", zap.Error(err))
		}
	}()

	var line []byte
	var client string

	// add this new client
	if client, err = s.addClient(conn); err != nil {
		s.Warn("error adding client.", zap.Error(err))
		return
	}

	// hand rest of the connection
	s.Info("handling connection ...",
		zap.String("remoteAddr", conn.RemoteAddr().String()))

	for {
		// receive request
		s.Info("waiting for request...", zap.String("client", client))
		reader := bufio.NewReader(conn)
		if line, _, err = reader.ReadLine(); err != nil {
			s.Warn("error reading from connection.", zap.Error(err))
			return
		}
		s.Info("received request",
			zap.String("request", string(line)),
			zap.String("client", client))

		data := make([]byte, len(line))
		copy(data, line)
		s.FanOut(client, append(data, '\n'))
	}
}

func (s *ChatServer) FanOut(sender string, line []byte) {
	wg := &sync.WaitGroup{}
	for clientName, clientConn := range s.clients {
		if clientName != sender {
			// LEARNING: this is a closure - ALWAYS pass the variables you need!
			wg.Add(1)
			go func(clientConn net.Conn, clientName string) {
				defer wg.Done()
				response := append([]byte(sender+": "), line...)
				if _, err := clientConn.Write(response); err != nil {
					s.Warn("failed to send message to client",
						zap.String("clientName", clientName),
						zap.Error(err))
				}
			}(clientConn, clientName)
		}
	}
	wg.Wait()
}

func (s *ChatServer) Listen() error {
	var err error
	var conn net.Conn
	var listener net.Listener

	s.Info("starting server...")
	if listener, err = net.Listen("tcp", ":7007"); err != nil {
		return fmt.Errorf("error listening: %v", err)
	}

	for {
		s.Info("waiting for connection...")
		if conn, err = listener.Accept(); err != nil {
			s.Info("error accepting connection.", zap.Error(err))
		}
		go func() {
			s.handleConnection(conn)
		}()
	}
}
