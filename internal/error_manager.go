package internal

import "sync"

type ErrorManager struct {
	mu *sync.Mutex
	e  *error
	es *[]error
}

func (em *ErrorManager) SetIfNotNil(err error) {
	if em == nil || err == nil {
		return
	}
	em.mu.Lock()
	defer em.mu.Unlock()
	if em.e != nil {
		*em.e = err
	}
	if em.es != nil {
		*em.es = append(*em.es, err)
	}
}

func (em *ErrorManager) ClearErr() {
	if em == nil || em.e == nil {
		return
	}
	em.mu.Lock()
	defer em.mu.Unlock()
	if em.e != nil {
		*em.e = nil
	}
}

func (em *ErrorManager) ClearErrs() {
	if em == nil || em.es == nil {
		return
	}
	em.mu.Lock()
	defer em.mu.Unlock()
	if em.es != nil {
		*em.es = nil
	}
}

func NewErrorManager(err *error, errs *[]error) *ErrorManager {
	if err == nil && errs == nil {
		return nil
	}
	return &ErrorManager{mu: &sync.Mutex{}, e: err, es: errs}
}
