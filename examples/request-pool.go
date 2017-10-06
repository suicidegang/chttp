package main

import (
	"log"
	"net/url"
)

func main() {
	start := time.Now()
	pool := Pool(100, 1000)

	for i := 0; i < 1000; i++ {

		// Send request signature into the pool to be managed.
		pool.In <- RequestOrPanic(GET("https://now.httpbin.org/"), Timeout(10))
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

	log.Printf("End processing requests in %v time.\n", time.Now().Sub(start))
}
