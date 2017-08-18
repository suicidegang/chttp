package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

type CHttp struct {
}

type SyncPool struct {
	Workers int
	In      chan<- chan Response
	Out     <-chan Response
}

func (pool SyncPool) Run(in chan chan Response, out chan Response) {
	for pid := 0; pid < pool.Workers; pid++ {
		go func(pid int) {
			for r := range in {
				log.Printf("STARTED #%v\n", pid)
				data := <-r
				out <- data
				log.Printf("OUT TO MAIN THREAD #%v.\n", pid)
			}
		}(pid)
	}
}

func Pool(n int) (pool SyncPool) {
	pool.Workers = n

	// Prepare channels and run workers.
	input := make(chan chan Response)
	output := make(chan Response, n)

	// Share channels...
	pool.In = input
	pool.Out = output
	pool.Run(input, output)
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pool := Pool(4)

	go func() {
		for i := 0; ; i++ {
			delay := rand.Intn(14)
			pool.In <- RequestOrPanic(GET("https://httpbin.org/delay/"+strconv.Itoa(delay)), Timeout(10))
		}
	}()

	log.Println("Sink consumer starting...")
	for r := range pool.Out {
		if r.Err != nil {
			log.Printf("Err while performing request %v", r.Err)
			continue
		}

		body, err := r.ReadAll()
		if err != nil {
			log.Printf("Err while reading data %v", err)
			continue
		}

		log.Printf("FINISHED MAIN THREAD, len %v...", len(body))
	}

}
