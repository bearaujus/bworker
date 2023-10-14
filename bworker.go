package bworker

import "sync"

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

func (bw *bWorker) Do(job Job) {
	if job == nil || bw.shutdown {
		return
	}
	job.queueToChan(bw.jobWG, bw.jobs)
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
	resetOptErrIfUsed(bw.mu, bw.optErr)
}

func (bw *bWorker) ResetErrs() {
	resetOptErrsIfUsed(bw.mu, bw.optErrs)
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
		job.do(bw.jobWG, bw.mu, bw.optRetry, bw.optErr, bw.optErrs)
	}
}
