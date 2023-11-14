package internal

import (
	"context"
	"sync"
)

type CtxManager struct {
	c    context.Context
	cl   context.CancelFunc
	rwMu *sync.RWMutex
}

func (cm *CtxManager) Ctx() context.Context {
	cm.rwMu.RLock()
	defer cm.rwMu.RUnlock()
	return cm.c
}

func (cm *CtxManager) Cancel() bool {
	cm.rwMu.Lock()
	defer cm.rwMu.Unlock()
	if cm.c.Err() != nil {
		return false
	}
	cm.cl()
	return true
}

func (cm *CtxManager) IsDead() bool {
	cm.rwMu.RLock()
	defer cm.rwMu.RUnlock()
	return cm.c.Err() != nil
}

func NewCtxManager() *CtxManager {
	c, cl := context.WithCancel(context.Background())
	return &CtxManager{
		c:    c,
		cl:   cl,
		rwMu: &sync.RWMutex{},
	}
}
