package bworker

import "sync"

// Job represent a function to be executed by a worker.
type Job func() error

func (j Job) execute(retry int, wg *sync.WaitGroup, mu *sync.Mutex, optErr *error, optErrs *[]error) {
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
		if optErr != nil {
			mu.Lock()
			*optErr = err
			mu.Unlock()
		}
		if optErrs != nil {
			mu.Lock()
			*optErrs = append(*optErrs, err)
			mu.Unlock()
		}
	}
}
