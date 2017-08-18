package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pool := Pool(15)

	for i := 0; i < 15; i++ {
		delay := rand.Intn(14)
		pool.In <- RequestOrPanic(GET("https://httpbin.org/delay/"+strconv.Itoa(delay)), Timeout(10))
	}

	// Emulate synchronous behavior by closing input channel.
	close(pool.In)

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

	log.Println("End processing requests.")
}
