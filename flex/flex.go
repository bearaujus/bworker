package flex

import (
	"github.com/bearaujus/bworker/internal"
)

type BWorkerFlex interface {
	// Do submit a job to be executed by a worker.
	Do(job func() error)

	// DoSimple submit a job to be executed by a worker without an error.
	DoSimple(job func())

	// Wait wait for all jobs to be completed.
	Wait()

	// ClearErr reset the error variable when you are using WithErrors.
	ClearErr()

	// ClearErrs reset the slice of error variables when you are using WithErrors.
	ClearErrs()
}

type bWorkerFlex struct {
	jobManager   *internal.JobManager
	errorManager *internal.ErrorManager
}

// NewBWorkerFlex create a new BWorkerFlex with OptionFlex(s) and unlimited concurrency level.
func NewBWorkerFlex(opts ...OptionFlex) BWorkerFlex {
	o := &internal.OptionFlex{}
	for _, opt := range opts {
		if opt != nil {
			opt.Apply(o)
		}
	}
	em := internal.NewErrorManager(o.Err, o.Errs)
	bwf := &bWorkerFlex{
		jobManager:   internal.NewJobManager(o.Retry, em),
		errorManager: em,
	}
	return bwf
}

func (bwf *bWorkerFlex) Do(job func() error) {
	if job == nil {
		return
	}
	pendingJob := bwf.jobManager.New(job)
	go pendingJob()
}

func (bwf *bWorkerFlex) DoSimple(job func()) {
	if job == nil {
		return
	}
	pendingJob := bwf.jobManager.NewSimple(job)
	go pendingJob()
}

func (bwf *bWorkerFlex) Wait() {
	bwf.jobManager.Wait()
}

func (bwf *bWorkerFlex) ClearErr() {
	bwf.errorManager.ClearErr()
}

func (bwf *bWorkerFlex) ClearErrs() {
	bwf.errorManager.ClearErrs()
}
