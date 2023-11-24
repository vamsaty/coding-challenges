package main

import (
	"coding-challenges/8-redis-server/server"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	s := server.NewRedisServer()
	PanicIf(s.Start())
}
