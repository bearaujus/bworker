package flex_option

import "github.com/bearaujus/bworker/internal"

type FlexOption interface {
	Apply(o *internal.FlexOption)
}

// WithRetry sets the number of times to retry a failed job.
func WithRetry(n int) FlexOption {
	return &withRetry{n}
}

type withRetry struct{ n int }

func (w *withRetry) Apply(o *internal.FlexOption) {
	if w.n <= 0 {
		return
	}
	o.Retry = w.n
}

// WithError sets a pointer to an error variable that will be populated if any job fails.
func WithError(err *error) FlexOption {
	return &withError{err}
}

type withError struct{ err *error }

func (w *withError) Apply(o *internal.FlexOption) {
	o.Err = internal.NewErrMem(w.err)
}

// WithErrors sets a pointer to a slice of error variables that will be populated if any jobs fail.
func WithErrors(errs *[]error) FlexOption {
	return &withErrors{errs}
}

type withErrors struct{ errs *[]error }

func (w *withErrors) Apply(o *internal.FlexOption) {
	o.Errs = internal.NewErrsMem(w.errs)
}
