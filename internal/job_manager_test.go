package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestJobManager(t *testing.T) {
	type args struct {
		numJobRetry int
		e           *error
		es          *[]error
	}
	tests := []struct {
		name        string
		args        args
		runner      func(jm *JobManager) *int64
		wantRet     int64
		wantErr     bool
		wantErrsLen int
	}{
		{
			name: "test use default value",
			args: args{
				numJobRetry: 0,
				e:           nil,
				es:          nil,
			},
			runner: func(jm *JobManager) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				j1 := jm.NewSimple(func() {
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				go j1()
				return &ret
			},
			wantRet:     1,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test without retry",
			args: args{
				numJobRetry: 0,
				e: func() *error {
					var err error
					return &err
				}(),
				es: func() *[]error {
					var errs []error
					return &errs
				}(),
			},
			runner: func(jm *JobManager) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				j1 := jm.NewSimple(func() {
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				go j1()
				j2 := jm.New(func() error {
					mu.Lock()
					defer mu.Unlock()
					ret++
					return errors.New("1")
				})
				go j2()
				return &ret
			},
			wantRet:     2,
			wantErr:     true,
			wantErrsLen: 1,
		},
		{
			name: "test with retry",
			args: args{
				numJobRetry: 10,
				e: func() *error {
					var err error
					return &err
				}(),
				es: func() *[]error {
					var errs []error
					return &errs
				}(),
			},
			runner: func(jm *JobManager) *int64 {
				var ret int64
				var mu = &sync.Mutex{}

				j1 := jm.NewSimple(func() {
					mu.Lock()
					defer mu.Unlock()
					ret++
				})
				go j1()
				j2 := jm.New(func() error {
					mu.Lock()
					defer mu.Unlock()
					ret++
					return errors.New("an error")
				})
				go j2()
				j3 := jm.New(func() error {
					mu.Lock()
					defer mu.Unlock()
					ret++
					return errors.New("an error")
				})
				go j3()
				return &ret
			},
			wantRet:     ((1 + 10) * 2) + 1, // ((base attempt + num retry)*num job with error)+do simple
			wantErr:     true,
			wantErrsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jm := NewJobManager(tt.args.numJobRetry, NewErrorManager(tt.args.e, tt.args.es))
			if tt.runner != nil {
				gotNumExecuted := tt.runner(jm)
				jm.Wait()
				assert.Equal(t, tt.wantRet, *gotNumExecuted)
			}
			if tt.args.e != nil {
				if tt.wantErr {
					assert.Error(t, *tt.args.e)
				} else {
					assert.NoError(t, *tt.args.e)
				}
			} else {
				assert.Nil(t, tt.args.e)
			}
			if tt.args.es != nil {
				assert.Equal(t, tt.wantErrsLen, len(*tt.args.es))
			} else {
				assert.Nil(t, tt.args.es)
			}
		})
	}
}
