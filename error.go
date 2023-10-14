package bworker

import "sync"

func setOptErrIfUsed(mu *sync.Mutex, optErr *error, err error) {
	if optErr == nil {
		return
	}

	mu.Lock()
	*optErr = err
	mu.Unlock()
}

func resetOptErrIfUsed(mu *sync.Mutex, optErr *error) {
	if optErr == nil {
		return
	}

	mu.Lock()
	*optErr = nil
	mu.Unlock()
}

func appendOptErrsIfUsed(mu *sync.Mutex, optErrs *[]error, err error) {
	if optErrs == nil {
		return
	}

	mu.Lock()
	*optErrs = append(*optErrs, err)
	mu.Unlock()
}

func resetOptErrsIfUsed(mu *sync.Mutex, optErrs *[]error) {
	if optErrs == nil {
		return
	}

	mu.Lock()
	*optErrs = nil
	mu.Unlock()
}
