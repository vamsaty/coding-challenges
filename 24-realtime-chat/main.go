package main

import (
	"os"
	"solid/substack-coding/24-realtime-chat/client"
	"solid/substack-coding/24-realtime-chat/server"
)

func main() {
	flag := os.Args[1]
	if flag == "server" {
		if err := server.NewChatServer().Listen(); err != nil {
			panic(err)
		}
	}
	if flag == "client" {
		if err := client.NewChatClient().Start(); err != nil {
			panic(err)
		}
	}
	if flag == "viewer" {
		if err := client.NewChatClient().StartViewer(); err != nil {
			panic(err)
		}
	}
}
