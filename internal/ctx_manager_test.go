package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCtxManager(t *testing.T) {
	tests := []struct {
		name   string
		runner func(cm *CtxManager)
	}{
		{
			name: "test basic use-cases",
			runner: func(cm *CtxManager) {
				assert.Nil(t, cm.Ctx().Err())
				assert.False(t, cm.IsDead())

				assert.True(t, cm.Cancel())

				assert.NotNil(t, cm.Ctx().Err())
				assert.True(t, cm.IsDead())

				assert.False(t, cm.Cancel())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCtxManager()
			assert.False(t, cm.IsDead())
			defer func() {
				cm.Cancel()
				assert.True(t, cm.IsDead())
			}()
			if tt.runner != nil {
				tt.runner(cm)
			}
		})
	}
}
