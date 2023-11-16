package main

import (
	"flag"
	"solid/open-source/coding-challenges/11-web-server/server"
)

var (
	httpHost = flag.String("host", "localhost", "hostname")
	httpPort = flag.Int("port", 9999, "port")
)

func main() {
	flag.Parse()

	config := &server.HttpServerConfig{
		Host: *httpHost,
		Port: *httpPort,
	}

	s := server.NewHttpServer(config)
	s.Register("GET", "/", server.FetchWebPage)
	s.StartServer()
}
