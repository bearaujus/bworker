package internal

import "sync"

type JobManager struct {
	wg      *sync.WaitGroup
	optErr  *ErrMem
	optErrs *ErrsMem
}

type PendingJob func(numRetry int)

func (jm *JobManager) New(job func() error) PendingJob {
	jm.wg.Add(1)
	return func(numRetry int) {
		defer jm.wg.Done()
		attempts := 1 + numRetry // 1 (base attempt) + num retry(s)
		for attempt := 0; attempt < attempts; attempt++ {
			err := job()
			if err == nil {
				return
			}
			if attempt != attempts-1 {
				continue
			}
			if jm.optErr != nil {
				jm.optErr.SetIfNotNil(err)
			}
			if jm.optErrs != nil {
				jm.optErrs.AppendIfNotNil(err)
			}
		}
	}
}

func (jm *JobManager) NewSimple(job func()) PendingJob {
	return jm.New(func() error {
		job()
		return nil
	})
}

func (jm *JobManager) Wait() {
	jm.wg.Wait()
}

func NewJobManager(optErr *ErrMem, optErrs *ErrsMem) *JobManager {
	return &JobManager{
		wg:      &sync.WaitGroup{},
		optErr:  optErr,
		optErrs: optErrs,
	}
}
