package bworker

type OptionFlex interface {
	Apply(flex *bWorkerFlex)
}

// WithRetryFlex sets the number of times to retry a failed job.
func WithRetryFlex(n int) OptionFlex {
	return &withRetryFlex{n}
}

type withRetryFlex struct{ n int }

func (w *withRetryFlex) Apply(bw *bWorkerFlex) {
	if w.n <= 0 {
		return
	}
	bw.optRetry = w.n
}

// WithErrorFlex sets a pointer to an error variable that will be populated if any job fails.
func WithErrorFlex(err *error) OptionFlex {
	return &withErrorFlex{err}
}

type withErrorFlex struct{ err *error }

func (w *withErrorFlex) Apply(bw *bWorkerFlex) {
	if w.err == nil {
		return
	}
	bw.optErr = w.err
}

// WithErrorsFlex sets a pointer to a slice of error variables that will be populated if any jobs fail.
func WithErrorsFlex(errs *[]error) OptionFlex {
	return &withErrorsFlex{errs}
}

type withErrorsFlex struct{ errs *[]error }

func (w *withErrorsFlex) Apply(bw *bWorkerFlex) {
	if w.errs == nil {
		return
	}
	bw.optErrs = w.errs
}
