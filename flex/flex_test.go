package flex

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestWorkerFlex(t *testing.T) {
	type args struct {
		opts []OptionFlex
	}
	tests := []struct {
		name        string
		args        args
		jobs        func(bwf BWorkerFlex) *atomic.Int64
		wantRet     int64
		wantErr     bool
		wantErrsLen int
	}{
		{
			name: "test use default value",
			args: args{
				opts: []OptionFlex{WithRetry(-1), nil},
			},
			jobs:        nil,
			wantRet:     0,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test execute nil job",
			args: args{
				opts: nil,
			},
			jobs: func(bwf BWorkerFlex) *atomic.Int64 {
				ret := &atomic.Int64{}

				bwf.Do(nil)
				bwf.DoSimple(nil)
				return ret
			},
			wantRet:     0,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test execute jobs with retry",
			args: args{
				opts: []OptionFlex{WithRetry(3), WithError(nil), WithErrors(nil)}, // Error will be masked at runner
			},
			jobs: func(bwf BWorkerFlex) *atomic.Int64 {
				ret := &atomic.Int64{}

				numJob, wantErrLen := 2000, 500
				for i := 0; i < numJob; i++ {
					icp := i
					bwf.Do(func() error {
						ret.Add(1)
						if icp < wantErrLen {
							return errors.New("an error")
						}
						return nil
					})
					bwf.DoSimple(func() { ret.Add(1) })

				}
				return ret
			},
			wantRet:     (500 * (1 + 3)) + 1500 + 2000, // (wantErrLen*(1+numRetry)) + doSuccessLen + doSimple
			wantErr:     true,
			wantErrsLen: 500,
		},
		{
			name: "test clear error",
			args: args{
				opts: []OptionFlex{WithError(nil), WithErrors(nil)},
			},
			jobs: func(bwf BWorkerFlex) *atomic.Int64 {
				ret := &atomic.Int64{}

				bwf.Do(func() error {
					ret.Add(1)
					return errors.New("an error")
				})
				bwf.Do(func() error {
					ret.Add(1)
					return errors.New("an error")
				})
				bwf.Wait()
				bwf.ClearErr()
				bwf.ClearErrs()
				return ret
			},
			wantRet:     2,
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
			bwp := NewBWorkerFlex(tt.args.opts...)
			if tt.jobs != nil {
				gotNumExecuted := tt.jobs(bwp)
				bwp.Wait()
				assert.Equal(t, tt.wantRet, gotNumExecuted.Load())
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
