package main

import (
	"context"
	"github.com/bean-du/pool"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init a pool with options
	client := pool.NewClient(
		pool.WebsocketDialer("ws://127.0.0.1:8081/ws"),
		// set pool size
		pool.WithPoolSize(50),
		// set write func default is tcp writer
		pool.WithWriteFunc(pool.WsWriter),
		// set read func, must be set
		pool.WithReadFunc(pool.WebsocketReadFunc(dataHandleFunc)),
		// set min idle connections
		pool.WithMinIdleConns(10),

		pool.WithPoolTimeout(5*time.Second),
		// set idle check duration
		pool.WithIdleCheckFrequency(time.Second*10),
	)

	for i := 0; i < 10; i++ {
		go func() {
			if err := client.Send(context.Background(), []byte("hello")); err != nil {
				log.Println(err)
			}
		}()
	}

	select {
	case <-sig:
		client.Close()
	}
}

func dataHandleFunc(p []byte) {
	go func() {
		log.Println(string(p))
	}()
}
