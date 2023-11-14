package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorManager(t *testing.T) {
	type args struct {
		e  *error
		es *[]error
	}
	tests := []struct {
		name        string
		args        args
		runner      func(em *ErrorManager)
		wantNil     bool
		wantErr     bool
		wantErrsLen int
	}{
		{
			name: "test use default value",
			args: args{
				e:  nil,
				es: nil,
			},
			runner:      nil,
			wantNil:     true,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test execute operation on nil instance",
			args: args{
				e:  nil,
				es: nil,
			},
			runner: func(em *ErrorManager) {
				em.SetIfNotNil(nil)
				em.SetIfNotNil(errors.New("an error"))
				em.ClearErr()
				em.ClearErrs()
			},
			wantNil:     true,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test set nil error",
			args: args{
				e: func() *error {
					var err error
					return &err
				}(),
				es: func() *[]error {
					var errs []error
					return &errs
				}(),
			},
			runner: func(em *ErrorManager) {
				em.SetIfNotNil(nil)
				em.SetIfNotNil(nil)
				em.SetIfNotNil(nil)
			},
			wantNil:     false,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test set not nil error",
			args: args{
				e: func() *error {
					var err error
					return &err
				}(),
				es: func() *[]error {
					var errs []error
					return &errs
				}(),
			},
			runner: func(em *ErrorManager) {
				em.SetIfNotNil(errors.New("an error"))
				em.SetIfNotNil(nil)
				em.SetIfNotNil(errors.New("an error"))
				em.SetIfNotNil(errors.New("an error"))
			},
			wantNil:     false,
			wantErr:     true,
			wantErrsLen: 3,
		},
		{
			name: "test clear",
			args: args{
				e: func() *error {
					var err error
					return &err
				}(),
				es: func() *[]error {
					var errs []error
					return &errs
				}(),
			},
			runner: func(em *ErrorManager) {
				em.SetIfNotNil(errors.New("an error"))
				em.ClearErr()
				em.ClearErrs()
			},
			wantNil:     false,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test clear err when not set",
			args: args{
				e: nil,
				es: func() *[]error {
					var errs []error
					return &errs
				}(),
			},
			runner: func(em *ErrorManager) {
				em.SetIfNotNil(errors.New("an error"))
				em.ClearErr()
				em.ClearErrs()
			},
			wantNil:     false,
			wantErr:     false,
			wantErrsLen: 0,
		},
		{
			name: "test clear errs when not set",
			args: args{
				e: func() *error {
					var err error
					return &err
				}(),
				es: nil,
			},
			runner: func(em *ErrorManager) {
				em.SetIfNotNil(errors.New("an error"))
				em.ClearErr()
				em.ClearErrs()
			},
			wantNil:     false,
			wantErr:     false,
			wantErrsLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := NewErrorManager(tt.args.e, tt.args.es)
			if tt.runner != nil {
				tt.runner(em)
			}
			if tt.wantNil {
				assert.Nil(t, em)
			} else {
				assert.NotNil(t, em)
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
