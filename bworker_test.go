package bworker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestBWorker(t *testing.T) {
	t.Run("use default num worker", func(t *testing.T) {
		bw := NewBWorker(0)
		defer bw.Shutdown()
	})

	t.Run("execute nil job", func(t *testing.T) {
		bw := NewBWorker(1)
		defer bw.Shutdown()

		bw.Do(nil)
	})

	t.Run("execute operation when already shut down", func(t *testing.T) {
		bw := NewBWorker(1)
		bw.Shutdown()

		bw.Do(func() error { return nil })
		bw.Wait()
		bw.Shutdown()
	})

	t.Run("use reset when option not set", func(t *testing.T) {
		bw := NewBWorker(1)
		defer bw.Shutdown()

		bw.ResetErr()
		bw.ResetErrs()
	})

	t.Run("test corner case for worker option", func(t *testing.T) {
		bw := NewBWorker(1,
			WithJobBuffer(0),
			WithRetry(0),
			WithError(nil),
			WithErrors(nil),
		)
		defer bw.Shutdown()
	})

	t.Run("basic use cases", func(t *testing.T) {
		var err error
		var errs []error
		numRetry := 3
		bw := NewBWorker(10,
			WithJobBuffer(10),
			WithRetry(numRetry),
			WithError(&err),
			WithErrors(&errs),
		)

		numJob, wantErrLen, retried, mu := 100, 50, 0, &sync.Mutex{}
		for i := 0; i < numJob; i++ {
			icp := i
			bw.Do(func() error {
				if icp < wantErrLen {
					mu.Lock()
					defer mu.Unlock()
					retried++
					return errors.New("an error")
				}
				return nil
			})
		}

		bw.Wait()
		assert.NotEmpty(t, errs)
		assert.Error(t, err)
		assert.Equal(t, wantErrLen, len(errs))
		for _, err := range errs {
			assert.Error(t, err)
		}
		assert.Equal(t, wantErrLen*(1+numRetry), retried)

		bw.ResetErr()
		bw.ResetErrs()

		assert.NoError(t, err)
		assert.Empty(t, errs)
	})
}
