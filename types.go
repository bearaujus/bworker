package bworker

type BWorker interface {
	// Do submit a job to be executed by a worker. If IsDead this function will perform no-op.
	// On worker type created by NewBWorker, this may block the thread.
	// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
	// Also, you can consider to use option.WithJobBuffer.
	Do(job func() error)

	// DoSimple submit a job to be executed by a worker without an error. If IsDead this function will perform no-op.
	DoSimple(job func())

	// Wait wait for all jobs to be completed. If IsDead this function will perform no-op.
	Wait()

	// Shutdown shut down the worker pool. After performing this operation, Do and DoSimple will perform no-op.
	// If IsDead this function will perform no-op.
	Shutdown()

	// IsDead indicates the BWorker is already shut down or not.
	IsDead() bool

	// ResetErr reset the error variable when you are using option.WithErrors or flex_option.WithErrors.
	ResetErr()

	// ResetErrs reset the slice of error variables when you are using option.WithErrors or flex_option.WithErrors.
	ResetErrs()
}
