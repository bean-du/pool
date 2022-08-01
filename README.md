## Connection Pool 

> Is a collection of connections to a client connection.


## Usage

```go
package main

import (
	"context"
	"github.com/bean-du/pool"
	"log"
	"os"
	"os/signal"
)

func main() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, os.Kill)
    
    // init a pool with default options
    client := pool.NewClient(
        pool.WebsocketDialer("ws://127.0.0.1:8081/ws"),
        pool.WithPoolSize(50),
        pool.WithPoolFIFO(true),
        pool.WithReadFunc(pool.WebsocketReadFunc(dataHandleFunc)),
        pool.WithMinIdleConns(5),
    )

    for i := 0; i < 10; i++ {
        if err := client.Send(context.Background(), []byte("hello")); err != nil {
         log.Println(err)
        }
    }

    select {
    case <-sig:
        client.Close()
    }
}

func dataHandleFunc(p []byte) {
    log.Println(string(p))
}
```