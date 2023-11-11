package bworker

import (
	"errors"
	"github.com/bearaujus/bworker/option"
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
		defer bw.Shutdown()

		assert.False(t, bw.IsDead())
		bw.Shutdown()
		assert.True(t, bw.IsDead())

		var doExecuted bool
		bw.Do(func() error {
			doExecuted = true
			return nil
		})

		var doSimpleExecuted bool
		bw.DoSimple(func() {
			doSimpleExecuted = true
		})

		bw.Wait()
		bw.Shutdown()

		assert.True(t, bw.IsDead())
		assert.False(t, doExecuted)
		assert.False(t, doSimpleExecuted)
	})

	t.Run("use reset when error option not set", func(t *testing.T) {
		bw := NewBWorker(1)
		defer bw.Shutdown()

		bw.ResetErr()
		bw.ResetErrs()
	})

	t.Run("occur error when error option not set", func(t *testing.T) {
		bw := NewBWorker(1)
		defer bw.Shutdown()

		bw.Do(func() error {
			return errors.New("an error")
		})
	})

	t.Run("test corner case for worker option", func(t *testing.T) {
		bw := NewBWorker(1,
			option.WithJobBuffer(0),
			option.WithRetry(0),
			option.WithError(nil),
			option.WithErrors(nil),
		)
		defer bw.Shutdown()
	})

	t.Run("do use cases", func(t *testing.T) {
		var err error
		var errs []error
		numRetry := 3

		bw := NewBWorker(10,
			option.WithJobBuffer(10),
			option.WithRetry(numRetry),
			option.WithError(&err),
			option.WithErrors(&errs),
		)
		defer bw.Shutdown()

		numJob, wantErrLen, retried, mu := 12345, 5678, 0, &sync.Mutex{}
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

	t.Run("do simple use cases", func(t *testing.T) {
		var err error
		var errs []error

		bw := NewBWorker(10,
			option.WithJobBuffer(10),
			option.WithRetry(3),
			option.WithError(&err),
			option.WithErrors(&errs),
		)
		defer bw.Shutdown()

		numJob, numExecuted, mu := 12345, 0, &sync.Mutex{}
		for i := 0; i < numJob; i++ {
			bw.DoSimple(func() {
				mu.Lock()
				defer mu.Unlock()
				numExecuted++
			})
		}

		bw.Wait()
		assert.Empty(t, errs)
		assert.NoError(t, err)
		assert.Equal(t, numJob, numExecuted)
	})
}
