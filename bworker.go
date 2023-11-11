package bworker

import (
	"github.com/bearaujus/bworker/internal"
	"github.com/bearaujus/bworker/option"
	"sync"
	"sync/atomic"
)

type bWorker struct {
	wgWorker   *sync.WaitGroup
	jobManager *internal.JobManager
	jobs       chan internal.PendingJob
	option     *internal.Option
	shutdown   *atomic.Bool
}

// NewBWorker creates a new worker pool with the specified concurrency level and Option(s).
func NewBWorker(concurrency int, opts ...option.Option) BWorker {
	// set to default concurrency
	if concurrency <= 0 {
		concurrency = 1
	}
	bwOpt := internal.Option{}
	for _, opt := range opts {
		opt.Apply(&bwOpt)
	}
	bw := bWorker{
		wgWorker:   &sync.WaitGroup{},
		jobManager: internal.NewJobManager(bwOpt.Err, bwOpt.Errs),
		jobs:       make(chan internal.PendingJob, bwOpt.JobBuffer), // expected if job buffer is 0
		option:     &bwOpt,
		shutdown:   &atomic.Bool{},
	}
	// create workers
	bw.wgWorker.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer bw.wgWorker.Done()
			for job := range bw.jobs {
				job(bw.option.Retry)
			}
		}()
	}
	return &bw
}

func (bw *bWorker) Do(job func() error) {
	if bw.shutdown.Load() || job == nil {
		return
	}
	pendingJob := bw.jobManager.New(job)
	bw.jobs <- pendingJob
}

func (bw *bWorker) DoSimple(job func()) {
	if bw.shutdown.Load() || job == nil {
		return
	}
	pendingJob := bw.jobManager.NewSimple(job)
	bw.jobs <- pendingJob
}

func (bw *bWorker) Wait() {
	if bw.shutdown.Load() {
		return
	}
	bw.jobManager.Wait()
}

func (bw *bWorker) Shutdown() {
	if !bw.shutdown.CompareAndSwap(false, true) {
		return
	}
	close(bw.jobs)
	bw.jobManager.Wait()
	bw.wgWorker.Wait()
}

func (bw *bWorker) IsDead() bool {
	return bw.shutdown.Load()
}

func (bw *bWorker) ResetErr() {
	if bw.option.Err == nil {
		return
	}
	bw.option.Err.Clear()
}

func (bw *bWorker) ResetErrs() {
	if bw.option.Errs == nil {
		return
	}
	bw.option.Errs.Clear()
}
