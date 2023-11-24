package client

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var messages = []string{
	"From fairest creatures we desire increase,",
	"That thereby beauty's rose might never die,",
	"But as the riper should by time decease,",
	"His tender heir might bear his memory:",
	"Capital is dead labour, which, vampire-like, lives only by sucking living labour, and lives the more, the more labour it sucks.",
	"The product of the capitalist, the capitalist himself, is labour-power.",
	"Labour-power is a commodity, not capital.",
	"Capital is not a thing, but a social relation between persons, established by the instrumentality of things.",
	"Capital is a social relation of production between capitalists and wage-labourers.",
	"Whoever fights monsters should see to it that in the process he does not become a monster.",
	"Hello, world!",
}

func TestNewChatClient(t *testing.T) {

	numClients := 100
	clients := make([]*ChatClient, numClients)
	start := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < numClients; i++ {
		clients[i] = NewChatClient()
		clients[i].name = fmt.Sprintf("client_%d", i)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			clients[i].connect()
		}(i)
	}
	wg.Wait()
	fmt.Printf("elapsed: %v\n", time.Since(start))

	wg = &sync.WaitGroup{}
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			clients[i].Start()
		}(i)
	}

}
