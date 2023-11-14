# BWorker - Simple & Efficient Worker Implementation in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bearaujus/bjson/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bearaujus/bworker)](https://goreportcard.com/report/github.com/bearaujus/bworker)

BWorker is a lightweight and easy-to-use Go library that provides a simple and efficient way to execute concurrent
tasks.

It is designed to be flexible enough to handle a wide range of use cases, from simple tasks such as initializing app
resources
to more complex tasks such as processing large use-cases.

## Installation

To install BWorker, you can run the following command:

```shell
go get github.com/bearaujus/bworker
```

## Usage Overview

Here is the simple usage overview for BWorker Pool:

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

Output:

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

## Worker Types

### 1. BWorker Pool

An BWorker instance with **specified** concurrency level.

import:

```go
import "github.com/bearaujus/bworker/pool"
```

Init:

```go
pool.NewBWorkerPool(concurrency int, opts ...OptionPool)
```

List available options:

>
> - **WithJobPoolSize** - Set the size of the job pool size.
> - **WithStartupStagger** - Set the worker pool to stagger the startup of workers with the calculated delay.
> - **WithRetry** - Set the number of times to retry a failed job.
> - **WithError** - Set a pointer to an error variable that will be populated if any job fails.
> - **WithErrors** - Set a pointer to a slice of error variables that will be populated if any job fails.

List available functions:

```go
// Do submit a job to be executed by a worker. If IsDead this function will perform no-op.
// This function may block the thread (see pool/pool_test.go for more details).
//
// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
// Also, you can consider using WithJobPoolSize.
Do(job func () error)

// DoSimple submit a job to be executed by a worker without an error. If IsDead this function will perform no-op.
// This function may block the thread (see pool/pool_test.go for more details).
//
// To avoid thread blocking, you can adjust the inputted job like adding context with deadline to it.
// Also, you can consider using WithJobPoolSize.
DoSimple(job func ())

// Wait wait for all jobPool to be completed. If IsDead this function will perform no-op.
Wait()

// Shutdown shut down the worker pool. After performing this operation, Do and DoSimple will perform no-op.
// If IsDead this function will perform no-op.
Shutdown()

// IsDead indicates the BWorkerPool is already shut down or not.
IsDead() bool

// ClearErr reset the error variable when you are using WithErrors.
ClearErr()

// ClearErrs reset the slice of error variables when you are using WithErrors.
ClearErrs()
```

### 2. BWorker Flex

An BWorker instance with **unlimited** concurrency level.

import:

```go
import "github.com/bearaujus/bworker/flex"
```

Init:

```go
flex.NewBWorkerFlex(opts ...OptionFlex)
```

List available options:

>
> - **WithRetry** - Set the number of times to retry a failed job.
> - **WithError** - Set a pointer to an error variable that will be populated if any job fails.
> - **WithErrors** - Set a pointer to a slice of error variables that will be populated if any job fails.

List available functions:

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

## TODO

- Improve documentation at [README.md](https://github.com/bearaujus/bworker/blob/master/README.md)
- Improve project description and tagging

## License

This project is licensed under the MIT License - see the LICENSE file for details.
