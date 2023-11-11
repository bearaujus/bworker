package internal

import "sync"

type ErrMem struct {
	mu *sync.Mutex
	e  *error
}

func (em *ErrMem) SetIfNotNil(err error) {
	if err == nil {
		return
	}
	em.mu.Lock()
	defer em.mu.Unlock()
	*em.e = err
}

func (em *ErrMem) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()
	*em.e = nil
}

func NewErrMem(e *error) *ErrMem {
	if e == nil {
		return nil
	}
	return &ErrMem{mu: &sync.Mutex{}, e: e}
}
