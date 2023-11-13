package pool

import (
	"github.com/bearaujus/bworker/internal"
	"time"
)

type OptionPool interface {
	Apply(o *internal.OptionPool)
}

// WithJobPoolSize set the size of the job pool size. If you're not using this option, the default job pool size is 1.
func WithJobPoolSize(n int) OptionPool {
	return &withJobPoolSize{n}
}

type withJobPoolSize struct{ n int }

func (w *withJobPoolSize) Apply(o *internal.OptionPool) {
	if w.n <= 0 {
		return
	}
	o.JobPoolSize = w.n
}

// WithWorkerStartupDelay set the delay interval between the startup of each worker in the worker pool.
//
// For example, if you specify 3 concurrencies and WithWorkerStartupDelay(time.Second),
// it will start worker 1 at 0ms, worker 2 at 500ms and worker 3 at 1000ms.
func WithWorkerStartupDelay(d time.Duration) OptionPool {
	return &withWorkerStartupDelay{d}
}

type withWorkerStartupDelay struct{ d time.Duration }

func (w *withWorkerStartupDelay) Apply(o *internal.OptionPool) {
	if w.d <= 0 {
		return
	}
	o.WorkerStartupDelay = w.d
}

// WithRetry set the number of times to retry a failed job.
func WithRetry(n int) OptionPool {
	return &withRetry{n}
}

type withRetry struct{ n int }

func (w *withRetry) Apply(o *internal.OptionPool) {
	if w.n <= 0 {
		return
	}
	o.Retry = w.n
}

// WithError set a pointer to an error variable that will be populated if any job fails.
func WithError(e *error) OptionPool {
	return &withError{e}
}

type withError struct{ e *error }

func (w *withError) Apply(o *internal.OptionPool) {
	o.Err = w.e
}

// WithErrors set a pointer to a slice of error variables that will be populated if any job fails.
func WithErrors(es *[]error) OptionPool {
	return &withErrors{es}
}

type withErrors struct{ es *[]error }

func (w *withErrors) Apply(o *internal.OptionPool) {
	o.Errs = w.es
}
