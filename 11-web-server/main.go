package main

import (
	"solid/open-source/coding-challenges/11-web-server/server"
)

func main() {
	config := &server.HttpServerConfig{Host: "localhost", Port: 9999}

	s := server.NewHttpServer(config)
	s.Register("GET", "/", server.FetchWebPage)
	s.StartServer()
}
