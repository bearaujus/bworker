package internal

import "sync"

type ErrsMem struct {
	mu *sync.Mutex
	es *[]error
}

func (esm *ErrsMem) AppendIfNotNil(err error) {
	if err == nil {
		return
	}
	esm.mu.Lock()
	defer esm.mu.Unlock()
	*esm.es = append(*esm.es, err)
}

func (esm *ErrsMem) Clear() {
	esm.mu.Lock()
	defer esm.mu.Unlock()
	*esm.es = nil
}

func NewErrsMem(es *[]error) *ErrsMem {
	if es == nil {
		return nil
	}
	return &ErrsMem{mu: &sync.Mutex{}, es: es}
}
