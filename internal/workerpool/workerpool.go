package workerpool

import (
	"sync"
)

type WorkDispatcher struct {
	input chan int
}

// New creates a new WorkDispatcher with worker pool
func New(numWorkers int, processor func(int) string) *WorkDispatcher {
	wd := &WorkDispatcher{
		input: make(chan int),
	}

	// Start workers
	workerChannels := make([]<-chan string, numWorkers)
	for i := range numWorkers {
		workerChannels[i] = wd.worker(processor)
	}

	// Start fan-in goroutine
	output := wd.fanIn(workerChannels...)

	// Start processing results
	go func() {
		for result := range output {
			println(result)
		}
	}()

	return wd
}

func (wd *WorkDispatcher) AddWork(item int) {
	wd.input <- item
}

func (wd *WorkDispatcher) Close() {
	close(wd.input)
}

func (wd *WorkDispatcher) worker(processor func(int) string) <-chan string {
	output := make(chan string)
	go func() {
		for n := range wd.input {
			output <- processor(n)
		}
		close(output)
	}()
	return output
}

func (wd *WorkDispatcher) fanIn(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	output := make(chan string)

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			for item := range c {
				output <- item
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
