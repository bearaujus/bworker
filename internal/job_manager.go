package internal

import "sync"

type JobManager struct {
	wg           *sync.WaitGroup
	numJobRetry  int
	errorManager *ErrorManager
}

type PendingJob func()

func (jm *JobManager) New(job func() error) PendingJob {
	jm.wg.Add(1)
	return func() {
		defer jm.wg.Done()
		attempts := 1 + jm.numJobRetry // 1 (base attempt) + num retry(s)
		for attempt := 0; attempt < attempts; attempt++ {
			err := job()
			if err == nil {
				return
			}
			if attempt != attempts-1 {
				continue
			}
			jm.errorManager.SetIfNotNil(err)
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

func NewJobManager(numJobRetry int, errorManager *ErrorManager) *JobManager {
	return &JobManager{
		wg:           &sync.WaitGroup{},
		numJobRetry:  numJobRetry,
		errorManager: errorManager,
	}
}
