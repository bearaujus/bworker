package bworker

import (
	"sync"
)

type BWorker interface {
	// Do submit a job to be executed by a worker.
	Do(job Job)

	// Wait for all jobs to be completed.
	Wait()

	// Shutdown shut down the worker pool.
	Shutdown()

	// ResetErr reset the error variable when you are using option WithError.
	ResetErr()

	// ResetErrs reset the slice of error variables when you are using option WithErrors.
	ResetErrs()
}

type bWorker struct {
	wg    *sync.WaitGroup
	mu    *sync.Mutex
	jobWG *sync.WaitGroup
	jobs  chan Job

	optJobBuffer int
	optRetry     int
	optErr       *error
	optErrs      *[]error

	shutdown bool
}

// Job represent a function to be executed by a worker.
type Job func() error

// NewBWorker creates a new worker pool with the specified concurrency level and Option(s).
func NewBWorker(concurrency int, opts ...Option) BWorker {
	bw := bWorker{
		wg:    &sync.WaitGroup{},
		mu:    &sync.Mutex{},
		jobWG: &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt.Apply(&bw)
	}
	bw.jobs = make(chan Job, bw.optJobBuffer)
	bw.startWorkers(concurrency)
	return &bw
}

func (bw *bWorker) startWorkers(numWorkers int) {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	bw.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go bw.startWorker()
	}
}

func (bw *bWorker) startWorker() {
	defer bw.wg.Done()
	for job := range bw.jobs {
		bw.execute(job)
	}
}

func (bw *bWorker) execute(job Job) {
	defer bw.jobWG.Done()
	attempts := 1 + bw.optRetry
	for attempt := 0; attempt < attempts; attempt++ {
		err := job()
		if err == nil {
			return
		}
		if attempt != attempts-1 {
			continue
		}
		if bw.optErr != nil {
			bw.mu.Lock()
			*bw.optErr = err
			bw.mu.Unlock()
		}
		if bw.optErrs != nil {
			bw.mu.Lock()
			*bw.optErrs = append(*bw.optErrs, err)
			bw.mu.Unlock()
		}
	}
}

func (bw *bWorker) Do(job Job) {
	if job == nil || bw.shutdown {
		return
	}
	bw.jobWG.Add(1)
	bw.jobs <- job
}

func (bw *bWorker) Wait() {
	if bw.shutdown {
		return
	}
	bw.jobWG.Wait()
}

func (bw *bWorker) Shutdown() {
	if bw.shutdown {
		return
	}
	bw.shutdown = true
	close(bw.jobs)
	bw.jobWG.Wait()
	bw.wg.Wait()
}

func (bw *bWorker) ResetErr() {
	if bw.optErr == nil {
		return
	}
	bw.mu.Lock()
	*bw.optErr = nil
	bw.mu.Unlock()
}

func (bw *bWorker) ResetErrs() {
	if bw.optErrs == nil {
		return
	}
	bw.mu.Lock()
	*bw.optErrs = nil
	bw.mu.Unlock()
}