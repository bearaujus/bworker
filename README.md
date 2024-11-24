# BWorker - Worker and Concurrency Implementation in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bearaujus/bworker/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bearaujus/bworker)](https://goreportcard.com/report/github.com/bearaujus/bworker)

This package introduces a lightweight and easy-to-use Go library that empowers you with a simple and efficient method to
execute concurrent tasks.

## Installation

To install BWorker, you can run the following command:

```shell
go get github.com/bearaujus/bworker
```

## Worker Types

### 1. BWorker Pool

An BWorker instance with **specified** concurrency level.

- Import:

```go
import "github.com/bearaujus/bworker/pool"
```

- Initialize:

```go
pool.NewBWorkerPool(concurrency int, opts ...OptionPool)
```

- List available options:

```go
// WithStartupStagger set the worker pool to stagger the startup of workers with the calculated delay.
//
// For example, if you set 3 concurrencies and 1s delay, it will start worker 1 at 0ms, worker 2 at 500ms,
// and worker 3 at 1000ms.
//
// This option will work if you set more than 1 concurrency since the first worker will always start immediately. Delay formula:
//	delay = d / time.Duration(concurrency-1)
func WithStartupStagger(d time.Duration) OptionPool

// WithRetry set the number of times to retry a failed job.
func WithRetry(n int) OptionPool

// WithError set a pointer to an error variable that will be populated if any job fails.
func WithError(e *error) OptionPool

// WithErrors set a pointer to a slice of error variables that will be populated if any job fails.
func WithErrors(es *[]error) OptionPool
```

- List available functions:

```go
// Do submit a job to be executed by a worker. If IsDead this function will perform no-op.
// This function may block the thread (see pool/pool_test.go for more details).
//
// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
// Also, you can consider using WithJobPoolSize.
func Do(job func () error)

// DoSimple submit a job to be executed by a worker without an error. If IsDead this function will perform no-op.
// This function may block the thread (see pool/pool_test.go for more details).
//
// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
// Also, you can consider using WithJobPoolSize.
func DoSimple(job func ())

// Wait wait for all jobPool to be completed. If IsDead this function will perform no-op.
func Wait()

// Shutdown shut down the worker pool. After performing this operation, Do and DoSimple will perform no-op.
// If IsDead this function will perform no-op.
func Shutdown()

// IsDead indicates the BWorkerPool is already shut down or not.
func IsDead() bool

// ClearErr reset the error variable when you are using WithErrors.
func ClearErr()

// ClearErrs reset the slice of error variables when you are using WithErrors.
func ClearErrs()
```

### 2. BWorker Flex

An BWorker instance with **unlimited** concurrency level.

- Import:

```go
import "github.com/bearaujus/bworker/flex"
```

- Initialize:

```go
flex.NewBWorkerFlex(opts ...OptionFlex)
```

- List available options:

```go
// WithRetry set the number of times to retry a failed job.
func WithRetry(n int) OptionFlex

// WithError set a pointer to an error variable that will be populated if any job fails.
func WithError(e *error) OptionFlex

// WithErrors set a pointer to a slice of error variables that will be populated if any job fails.
func WithErrors(es *[]error) OptionFlex
```

- List available functions:

```go
// Do submit a job to be executed by a worker.
Do(job func () error)

// DoSimple submit a job to be executed by a worker without an error.
DoSimple(job func ())

// Wait wait for all jobs to be completed.
Wait()

// ClearErr reset the error variable when you are using WithErrors.
ClearErr()

// ClearErrs reset the slice of error variables when you are using WithErrors.
ClearErrs()
```

## Usage Example

```go
package main

import (
	"fmt"
	"github.com/bearaujus/bworker/pool"
)

func main() {
	// Init error vars
	var (
		err  error
		errs []error
	)

	// Init BWorker with some option(s)
	w := pool.NewBWorkerPool(
		// Num workers (concurrencies level)
		2,
		// Set the number of times to retry a failed job. Will executed 3 times on a failed job: 1 (base attempt) + 2 (retry)
		pool.WithRetry(2),
		// Set a pointer to an error variable that will be populated if any job fails
		pool.WithError(&err),
		// Set a pointer to a slice of error variables that will be populated if any job fails
		pool.WithErrors(&errs),
	)

	// Don't forget to add this block to avoid memory leak from the unclosed channel(s)
	defer w.Shutdown()

	// Add 2 jobs
	for i := 1; i <= 2; i++ {
		icp := i
		w.Do(func() error {
			fmt.Printf("job %v retrying...\n", icp)
			return fmt.Errorf("job %v error", icp)
		})
	}

	// Wait until all jobs executed
	w.Wait()

	fmt.Println(err)
	fmt.Println(errs)
}
```

```go
job 1 retrying...
job 1 retrying...
job 1 retrying...
job 2 retrying...
job 2 retrying...
job 2 retrying...
job 2 error
[job 1 error job 2 error]
```

## TODO

- Improve documentation at [README.md](https://github.com/bearaujus/bworker/blob/master/README.md)
- Improve project description and tagging

## License

This project is licensed under the MIT License - see
the [LICENSE](https://github.com/bearaujus/bworker/blob/master/LICENSE) file for details.
