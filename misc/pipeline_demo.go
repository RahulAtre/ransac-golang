/*
* Go Concurrency Pipeline example from https://go.dev/blog/pipelines
*/
package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <- done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(cs))
	out := make(chan int, 1)

	// new gorountine for for input channel to copy val to output chan
	output := func(c <-chan int) {
		// ensure done is called on return path for each output go rountine
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <- done:
				return
			}
		}
	}
	for _, c := range cs {
		go output(c)
	}

	// separate go routine which closes out channel after all input channels are closed
	// synchronized using WaitGroup
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// main function
func main()  {
	// set up pipeline
	// gen func provides the outbound channel
	c1 := gen(5,10)
	c2 := gen(2,4)

	// cancellation
	done := make(chan struct{})
	defer close(done)
	
	for n := range merge(done, sq(done, c1), sq(done, c2)) {
		fmt.Println(n)
	}
}
