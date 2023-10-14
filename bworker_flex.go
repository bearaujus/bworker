package bworker

import (
	"sync"
)

type bWorkerFlex struct {
	mu    *sync.Mutex
	jobWG *sync.WaitGroup

	optRetry int
	optErr   *error
	optErrs  *[]error

	shutdown bool
}

// NewBWorkerFlex creates a new worker pool with the Option(s) and unlimited concurrency level.
func NewBWorkerFlex(opts ...OptionFlex) BWorker {
	bwf := bWorkerFlex{
		mu:    &sync.Mutex{},
		jobWG: &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt.Apply(&bwf)
	}
	return &bwf
}

func (bwf *bWorkerFlex) Do(job Job) {
	if job == nil || bwf.shutdown {
		return
	}
	bwf.jobWG.Add(1)
	go job.execute(bwf.optRetry, bwf.jobWG, bwf.mu, bwf.optErr, bwf.optErrs)
}

func (bwf *bWorkerFlex) Wait() {
	if bwf.shutdown {
		return
	}
	bwf.jobWG.Wait()
}

func (bwf *bWorkerFlex) Shutdown() {
	if bwf.shutdown {
		return
	}
	bwf.shutdown = true
	bwf.jobWG.Wait()
}

func (bwf *bWorkerFlex) ResetErr() {
	if bwf.optErr == nil {
		return
	}
	bwf.mu.Lock()
	*bwf.optErr = nil
	bwf.mu.Unlock()
}

func (bwf *bWorkerFlex) ResetErrs() {
	if bwf.optErrs == nil {
		return
	}
	bwf.mu.Lock()
	*bwf.optErrs = nil
	bwf.mu.Unlock()
}
