package main

import (
	"sync"
)

type SyncPool struct {
	Workers int
	In      chan<- chan Response
	Out     <-chan Response
}

func (pool SyncPool) Run(in chan chan Response, out chan Response) {
	var done sync.WaitGroup

	// Done waiting group.
	done.Add(pool.Workers)

	for pid := 0; pid < pool.Workers; pid++ {

		// Starts a new worker that will range over input channel
		// waits for the response to come in
		// and finally sends it back to output channel. SRPful
		go func(pid int) {
			defer done.Done()

			for r := range in {
				data := <-r
				out <- data
			}
		}(pid)
	}

	// A separate process will wait for stuff to be done
	// and eventually (may or may not happen) close the output channel
	go func() {
		done.Wait()
		close(out)
	}()
}

// Creates a requests SyncPool with n workers.
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
