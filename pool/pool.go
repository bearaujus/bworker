package pool

import (
	"github.com/bearaujus/bworker/internal"
	"sync"
	"time"
)

type BWorkerPool interface {
	// Do submit a job to be executed by a worker. If IsDead this function will perform no-op.
	// This function may block the thread (see pool/pool_test.go for more details).
	//
	// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
	// Also, you can consider using WithJobPoolSize.
	Do(job func() error)

	// DoSimple submit a job to be executed by a worker without an error. If IsDead this function will perform no-op.
	// This function may block the thread (see pool/pool_test.go for more details).
	//
	// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
	// Also, you can consider using WithJobPoolSize.
	DoSimple(job func())

	// Wait wait for all jobPool to be completed. If IsDead this function will perform no-op.
	Wait()

	// Shutdown shut down the worker pool. After performing this operation, Do and DoSimple will perform no-op.
	// If IsDead this function will perform no-op.
	Shutdown()

	// IsDead indicates the BWorkerPool is already shut down or not.
	IsDead() bool

	// ClearErr reset the error variable when you are using WithErrors.
	ClearErr()

	// ClearErrs reset the slice of error variables when you are using WithErrors.
	ClearErrs()
}

type bWorkerPool struct {
	ctxManager   *internal.CtxManager
	jobManager   *internal.JobManager
	jobPool      chan internal.PendingJob
	errorManager *internal.ErrorManager
	wgWorker     *sync.WaitGroup
}

// NewBWorkerPool create a new BWorkerPool with OptionPool(s) and specified concurrency level.
//
// Please use BWorkerPool.Shutdown() to avoid memory leak from the unclosed channel(s).
func NewBWorkerPool(concurrency int, opts ...OptionPool) BWorkerPool {
	if concurrency <= 0 {
		concurrency = 1
	}
	o := &internal.OptionPool{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.Apply(o)
	}
	em := internal.NewErrorManager(o.Err, o.Errs)
	bwp := &bWorkerPool{
		ctxManager: internal.NewCtxManager(),
		jobManager: internal.NewJobManager(o.Retry, em),
		// If o.JobPoolSize = 0. It's basically the same with o.JobPoolSize = 1
		jobPool:      make(chan internal.PendingJob, o.JobPoolSize),
		errorManager: em,
		wgWorker:     &sync.WaitGroup{},
	}
	var startupDelay time.Duration
	if concurrency != 1 && o.StartupStagger != 0 {
		startupDelay = o.StartupStagger / time.Duration(concurrency-1)
	}
	bwp.wgWorker.Add(concurrency)
	go func() {
		for i := 0; i < concurrency; i++ {
			// the first worker will always start, before using startupDelay when using WithStartupStagger
			if i != 0 && o.StartupStagger != 0 {
				select {
				case <-time.Tick(startupDelay):
				case <-bwp.ctxManager.Ctx().Done():
				}
			}
			// Create a worker
			go func() {
				defer bwp.wgWorker.Done()
				// Keep pulling jobs until bw.jobPool is closed
				for job := range bwp.jobPool {
					job()
				}
			}()
		}
	}()
	return bwp
}

func (bwp *bWorkerPool) Do(job func() error) {
	if bwp.ctxManager.IsDead() || job == nil {
		return
	}
	pendingJob := bwp.jobManager.NewJob(job)
	bwp.jobPool <- pendingJob
}

func (bwp *bWorkerPool) DoSimple(job func()) {
	if bwp.ctxManager.IsDead() || job == nil {
		return
	}
	pendingJob := bwp.jobManager.NewJobSimple(job)
	bwp.jobPool <- pendingJob
}

func (bwp *bWorkerPool) Wait() {
	if bwp.ctxManager.IsDead() {
		return
	}
	// Wait until all jobs executed
	bwp.jobManager.Wait()
}

func (bwp *bWorkerPool) Shutdown() {
	if !bwp.ctxManager.Cancel() {
		return
	}
	// Wait until all jobs executed
	bwp.jobManager.Wait()
	// Shut down all active workers
	close(bwp.jobPool)
	// Wait until all workers are dead
	bwp.wgWorker.Wait()
}

func (bwp *bWorkerPool) IsDead() bool {
	return bwp.ctxManager.IsDead()
}

func (bwp *bWorkerPool) ClearErr() {
	bwp.errorManager.ClearErr()
}

func (bwp *bWorkerPool) ClearErrs() {
	bwp.errorManager.ClearErrs()
}
