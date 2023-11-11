package bworker

import (
	"github.com/bearaujus/bworker/flex_option"
	"github.com/bearaujus/bworker/internal"
	"sync/atomic"
)

type bWorkerFlex struct {
	jobManager *internal.JobManager
	option     *internal.FlexOption
	shutdown   *atomic.Bool
}

// NewFlexBWorker creates a new worker pool with the FlexOption(s) and unlimited concurrency level.
func NewFlexBWorker(opts ...flex_option.FlexOption) BWorker {
	bwOptFlex := internal.FlexOption{}
	for _, opt := range opts {
		opt.Apply(&bwOptFlex)
	}
	bwf := bWorkerFlex{
		jobManager: internal.NewJobManager(bwOptFlex.Err, bwOptFlex.Errs),
		option:     &bwOptFlex,
		shutdown:   &atomic.Bool{},
	}
	return &bwf
}

func (bwf *bWorkerFlex) Do(job func() error) {
	if bwf.shutdown.Load() || job == nil {
		return
	}
	pendingJob := bwf.jobManager.New(job)
	go pendingJob(bwf.option.Retry)
}

func (bwf *bWorkerFlex) DoSimple(job func()) {
	if bwf.shutdown.Load() || job == nil {
		return
	}
	pendingJob := bwf.jobManager.NewSimple(job)
	go pendingJob(bwf.option.Retry)
}

func (bwf *bWorkerFlex) Wait() {
	if bwf.shutdown.Load() {
		return
	}
	bwf.jobManager.Wait()
}

func (bwf *bWorkerFlex) Shutdown() {
	if !bwf.shutdown.CompareAndSwap(false, true) {
		return
	}
	bwf.jobManager.Wait()
}

func (bwf *bWorkerFlex) IsDead() bool {
	return bwf.shutdown.Load()
}

func (bwf *bWorkerFlex) ResetErr() {
	if bwf.option.Err == nil {
		return
	}
	bwf.option.Err.Clear()
}

func (bwf *bWorkerFlex) ResetErrs() {
	if bwf.option.Errs == nil {
		return
	}
	bwf.option.Errs.Clear()
}
