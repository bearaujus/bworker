package internal

type Option struct {
	JobBuffer int
	Retry     int
	Err       *ErrMem
	Errs      *ErrsMem
}

type FlexOption struct {
	Retry int
	Err   *ErrMem
	Errs  *ErrsMem
}
