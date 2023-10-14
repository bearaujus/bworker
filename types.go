package bworker

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
