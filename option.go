package bworker

type Option interface {
	Apply(*bWorker)
}

// WithJobBuffer sets the size of the job buffer.
func WithJobBuffer(n int) Option {
	return &withJobBuffer{n}
}

type withJobBuffer struct{ n int }

func (w *withJobBuffer) Apply(bw *bWorker) {
	if w.n <= 0 {
		return
	}
	bw.optJobBuffer = w.n
}

// WithRetry sets the number of times to retry a failed job.
func WithRetry(n int) Option {
	return &withRetry{n}
}

type withRetry struct{ n int }

func (w *withRetry) Apply(bw *bWorker) {
	if w.n <= 0 {
		return
	}
	bw.optRetry = w.n
}

// WithError sets a pointer to an error variable that will be populated if any job fails.
func WithError(err *error) Option {
	return &withError{err}
}

type withError struct{ err *error }

func (w *withError) Apply(bw *bWorker) {
	if w.err == nil {
		return
	}
	bw.optErr = w.err
}

// WithErrors sets a pointer to a slice of error variables that will be populated if any jobs fail.
func WithErrors(errs *[]error) Option {
	return &withErrors{errs}
}

type withErrors struct{ errs *[]error }

func (w *withErrors) Apply(bw *bWorker) {
	if w.errs == nil {
		return
	}
	bw.optErrs = w.errs
}
