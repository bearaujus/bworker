package pool

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	type args struct {
		concurrency int
		opts        []OptionPool
	}
	tests := []struct {
		name        string
		args        args
		jobs        func(bwp BWorkerPool) *int64
		wantRet     int64
		wantErr     bool
		wantErrsLen int
	}{
		{
			name: "test use default value",
			args: args{
				concurrency: -1,
				opts:        []OptionPool{WithJobPoolSize(-1), WithStartupStagger(-1), WithRetry(-1), nil},
			},
			jobs:        nil,
			wantRet:     0,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test execute nil job",
			args: args{
				concurrency: 50,
				opts:        nil,
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64

				bwp.Do(nil)
				bwp.DoSimple(nil)
				return &ret
			},
			wantRet:     0,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test execute jobs when already shut down",
			args: args{
				concurrency: 50,
				opts:        nil,
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				bwp.Shutdown()
				bwp.Do(func() error {
					mu.Lock()
					defer mu.Unlock()
					ret++
					return nil
				})
				bwp.DoSimple(func() {
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				return &ret
			},
			wantRet:     0,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test execute jobs with retry",
			args: args{
				concurrency: 50,
				opts:        []OptionPool{WithRetry(3), WithError(nil), WithErrors(nil)}, // Error will be masked at runner
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				numJob, wantErrLen := 2000, 500
				for i := 0; i < numJob; i++ {
					icp := i
					bwp.Do(func() error {
						mu.Lock()
						defer mu.Unlock()
						ret++
						if icp < wantErrLen {
							return errors.New("an error")
						}
						return nil
					})
					bwp.DoSimple(func() {
						mu.Lock()
						defer mu.Unlock()
						ret++
					})
				}
				return &ret
			},
			wantRet:     (500 * (1 + 3)) + 1500 + 2000, // (wantErrLen*(1+numRetry)) + doSuccessLen + doSimple
			wantErr:     true,
			wantErrsLen: 500,
		},
		{
			name: "test worker startup delay",
			args: args{
				concurrency: 2,
				opts:        []OptionPool{WithStartupStagger(time.Second)},
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				start := time.Now()
				bwp.DoSimple(func() { // Executed by w1
					time.Sleep(time.Second)
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				// After 1 sec
				bwp.DoSimple(func() { // Executed by w1
					time.Sleep(time.Second)
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				bwp.DoSimple(func() { // Executed by w2
					time.Sleep(time.Second)
					mu.Lock()
					defer mu.Unlock()
					ret++
				})

				bwp.Wait()
				// The total executed time should be around ~2 secs
				ts := time.Since(start)
				assert.LessOrEqual(t, time.Second*2, ts)
				assert.LessOrEqual(t, ts, (time.Second*2)+(time.Millisecond*100)) // Add 0.1s as a threshold
				return &ret
			},
			wantRet:     3,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test worker startup delay shutdown in the middle",
			args: args{
				concurrency: 200,
				opts:        []OptionPool{WithStartupStagger(time.Hour)},
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64

				start := time.Now()
				bwp.Shutdown()
				// The total block time should 0
				ts := time.Since(start)
				assert.LessOrEqual(t, ts, time.Millisecond*100) // Add 0.1s as a threshold
				return &ret
			},
			wantRet:     0,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test clear error",
			args: args{
				concurrency: 50,
				opts:        []OptionPool{WithError(nil), WithErrors(nil)},
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				bwp.Do(func() error {
					mu.Lock()
					defer mu.Unlock()
					ret++
					return errors.New("an error")
				})
				bwp.Do(func() error {
					mu.Lock()
					defer mu.Unlock()
					ret++
					return errors.New("an error")
				})
				bwp.Wait()
				bwp.ClearErr()
				bwp.ClearErrs()
				return &ret
			},
			wantRet:     2,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test blocked with jobs pool size",
			args: args{
				concurrency: 1,
				opts:        []OptionPool{WithJobPoolSize(1)},
			},
			jobs: func(bwp BWorkerPool) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				start := time.Now()
				bwp.DoSimple(func() { // Consumed by the worker (not blocking)
					time.Sleep(time.Second)
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				bwp.DoSimple(func() { // Queued at pool (not blocking)
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				bwp.DoSimple(func() { // Blocked (blocking)
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				// The total block time should be around ~ 1sec
				ts := time.Since(start)
				assert.LessOrEqual(t, time.Second, ts)
				assert.LessOrEqual(t, ts, (time.Second)+(time.Millisecond*100)) // Add 0.1s as a threshold
				return &ret
			},
			wantRet:     3,
			wantErr:     false,
			wantErrsLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err  error
				errs []error
			)
			for _, opt := range tt.args.opts {
				switch o := opt.(type) {
				case *withError:
					o.e = &err
				case *withErrors:
					o.es = &errs
				}
			}
			bwp := NewBWorkerPool(tt.args.concurrency, tt.args.opts...)
			assert.False(t, bwp.IsDead())
			defer func() {
				bwp.Shutdown()
				assert.True(t, bwp.IsDead())
			}()
			if tt.jobs != nil {
				gotNumExecuted := tt.jobs(bwp)
				bwp.Wait()
				assert.Equal(t, tt.wantRet, *gotNumExecuted)
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantErrsLen, len(errs))
			for _, v := range errs {
				assert.Error(t, v)
			}
		})
	}
}
