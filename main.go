package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Two workers, 14 buffer capacity.
	pool := Pool(2, 14)

	for i := 0; i < 15; i++ {
		delay := rand.Intn(14)

		// Send request signature into the pool to be managed.
		pool.In <- RequestOrPanic(GET("https://httpbin.org/delay/"+strconv.Itoa(delay)), Timeout(10))
	}

	// Emulate synchronous behavior by closing input channel.
	close(pool.In)

	log.Println("Sink consumer starting...")
	for response := range pool.Out {
		if response.Err != nil {
			log.Printf("Err while performing request %v", response.Err)
			continue
		}

		body, err := response.ReadAll()
		if err != nil {
			log.Printf("Err while reading data %v", err)
			continue
		}

		log.Printf("FINISHED MAIN THREAD, len %v...", len(body))
	}

	log.Println("End processing requests.")
}
