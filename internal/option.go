package internal

import "time"

type OptionPool struct {
	JobPoolSize    int
	StartupStagger time.Duration
	Retry          int
	Err            *error
	Errs           *[]error
}

type OptionFlex struct {
	Retry int
	Err   *error
	Errs  *[]error
}
