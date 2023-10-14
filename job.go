package bworker

import "sync"

// Job represent a function to be executed by a worker.
type Job func() error

func (j Job) executeInBackground(wg *sync.WaitGroup, mu *sync.Mutex, retry int, optErr *error, optErrs *[]error) {
	wg.Add(1)
	go j.do(wg, mu, retry, optErr, optErrs)
}

func (j Job) queueToChan(wg *sync.WaitGroup, c chan Job) {
	wg.Add(1)
	c <- j
}

func (j Job) do(wg *sync.WaitGroup, mu *sync.Mutex, retry int, optErr *error, optErrs *[]error) {
	defer wg.Done()
	attempts := 1 + retry
	for attempt := 0; attempt < attempts; attempt++ {
		err := j()
		if err == nil {
			return
		}
		if attempt != attempts-1 {
			continue
		}
		setOptErrIfUsed(mu, optErr, err)
		appendOptErrsIfUsed(mu, optErrs, err)
	}
}
