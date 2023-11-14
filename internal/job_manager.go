package internal

import "sync"

type JobManager struct {
	wg  *sync.WaitGroup
	njr int
	em  *ErrorManager
}

type PendingJob func()

func (jm *JobManager) NewJob(job func() error) PendingJob {
	jm.wg.Add(1)
	return func() {
		defer jm.wg.Done()
		ats := 1 + jm.njr // 1 (base attempt) + num retry(s)
		for at := 0; at < ats; at++ {
			err := job()
			if err == nil {
				return
			}
			if at != ats-1 {
				continue
			}
			jm.em.SetIfNotNil(err)
		}
	}
}

func (jm *JobManager) NewJobSimple(job func()) PendingJob {
	return jm.NewJob(func() error {
		job()
		return nil
	})
}

func (jm *JobManager) Wait() {
	jm.wg.Wait()
}

func NewJobManager(numJobRetry int, errorManager *ErrorManager) *JobManager {
	return &JobManager{
		wg:  &sync.WaitGroup{},
		njr: numJobRetry,
		em:  errorManager,
	}
}
