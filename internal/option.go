package internal

import "time"

type OptionPool struct {
	JobPoolSize        int
	WorkerStartupDelay time.Duration
	Retry              int
	Err                *error
	Errs               *[]error
}

type OptionFlex struct {
	Retry int
	Err   *error
	Errs  *[]error
}
