package flex

import "github.com/bearaujus/bworker/internal"

type OptionFlex interface {
	Apply(o *internal.OptionFlex)
}

// WithRetry set the number of times to retry a failed job.
func WithRetry(n int) OptionFlex {
	return &withRetry{n}
}

type withRetry struct{ n int }

func (w *withRetry) Apply(o *internal.OptionFlex) {
	if w.n <= 0 {
		return
	}
	o.Retry = w.n
}

// WithError set a pointer to an error variable that will be populated if any job fails.
func WithError(e *error) OptionFlex {
	return &withError{e}
}

type withError struct{ e *error }

func (w *withError) Apply(o *internal.OptionFlex) {
	o.Err = w.e
}

// WithErrors set a pointer to a slice of error variables that will be populated if any job fails.
func WithErrors(es *[]error) OptionFlex {
	return &withErrors{es}
}

type withErrors struct{ es *[]error }

func (w *withErrors) Apply(o *internal.OptionFlex) {
	o.Errs = w.es
}
