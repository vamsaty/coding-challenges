package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type ChatClient struct {
	name        string
	conn        net.Conn
	inputReader *bufio.Reader
	exitChan    chan bool
}

func NewChatClient() *ChatClient {
	return &ChatClient{
		inputReader: bufio.NewReader(os.Stdin),
		exitChan:    make(chan bool, 1),
	}
}

func (c *ChatClient) inputName() {
	fmt.Printf("name: ")
	line, _, err := c.inputReader.ReadLine()
	if err != nil {
		panic(err)
	}
	c.name = string(line)
}

func (c *ChatClient) connect() {
	var err error
	c.conn, err = net.Dial("tcp", "127.0.0.1:7007")
	if err != nil {
		panic(err)
	}

	_, err = c.conn.Write([]byte(c.name + "\n"))
	if err != nil {
		panic(err)
	}
}

func (c *ChatClient) Start() (err error) {
	c.inputName()
	c.connect()

	fmt.Printf("welcome to the chatroom : %s\n", c.name)
	go c.awaitRequest()
	go c.awaitResponse()

	<-c.exitChan
	return err
}

func (c *ChatClient) awaitResponse() {
	var line []byte
	var err error

	for {
		// receive response
		if line, _, err = bufio.NewReader(c.conn).ReadLine(); err == io.EOF {
			c.exitChan <- true
			return
		} else if err != nil {
			fmt.Printf("error : %v\n", err)
			continue
		}
		fmt.Printf("%v\n", string(line))
	}
}

func (c *ChatClient) awaitRequest() {
	var line []byte
	var err error

	for {
		// read user input
		if line, _, err = c.inputReader.ReadLine(); err == io.EOF {
			c.exitChan <- true
			return
		} else if err != nil {
			fmt.Printf("error : %v\n", err)
			continue
		} else if len(line) == 4 && string(line) == "quit" || string(line) == "exit" {
			c.exitChan <- true
			return
		}
		line = append(line, '\n')

		// send request
		if _, err = c.conn.Write(line); err == io.EOF {
			fmt.Println("failed to send request to server")
			c.exitChan <- true
			return
		}
	}
}
