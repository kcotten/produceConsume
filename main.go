package main

import (
    "fmt"
    "time"
)

func nextInt(prev int) int {
    time.Sleep(1 * time.Second) // pretend this is an expensive operation
    return prev + 1
}

func main() {
    prod := &producer{
        data: make(chan int),
        quit: make(chan chan error),
    }

    // producer
    go func() {
        var i = 0
        for {
            i = nextInt(i)
            select {
            case prod.data <- i:
            case ch := <-prod.quit:
                close(prod.data)
                // If the producer had an error while shutting down,
                // we could write the error to the ch channel here.
                close(ch)
                return
            }
        }
    }()

    // consumer
    for i := range prod.data {
        fmt.Printf("i=%v\n", i)
        if i >= 5 {
            err := prod.CloseHandler()
            if err != nil {
                // cannot happen in this example
                fmt.Printf("unexpected error: %v\n", err)
            }
        }
    }
}

type producer struct {
    data chan int
    quit chan chan error
}

// Handle exit signal
func (p *producer) CloseHandler() error {
	ch := make(chan error)
    p.quit <- ch
    return <-ch
}
