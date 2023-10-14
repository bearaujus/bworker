# BWorker - Simple & Efficient Worker Pool Implementation in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bearaujus/bjson/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bearaujus/bworker)](https://goreportcard.com/report/github.com/bearaujus/bworker)

BWorker is a lightweight and easy-to-use Go library that provides a simple and efficient way to execute concurrent
tasks. It is designed to be flexible enough to handle a wide range of use cases, from simple tasks such as sending
emails to more complex tasks such as processing large datasets.

## Features

- **Concurrency:** BWorker provides a simple way to control the number of concurrent workers in your pool.
- **Job queue:** BWorker allows you to queue jobs to be executed by your workers, and it will automatically manage the
  execution of those jobs.
- **Error handling:** BWorker provides a way to handle errors that occur during job execution.
- **Job retry**: BWorker can automatically retry jobs that fail, a specified number of times.
- **Shutdown:** BWorker provides a way to gracefully shut down your worker pool.

## Installation

To install BWorker, you can run the following command:

```shell
go get github.com/bearaujus/bjson
```

## Usage

To use BWorker, simply create a new instance of the `BWorker` and specify the number of concurrent workers you want to
use. You can then add jobs to the worker pool using the `Do()` method. The worker pool will automatically start
processing the jobs in the queue.

When you are finished with your worker pool, you should call the `Shutdown()` method to gracefully shut it down.

### Examples

The following example shows how to use BWorker to execute a simple job:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/bearaujus/bworker"
)

func main() {
	// Create a new BWorker instance with 10 concurrent workers and a retry count of 2.
	bw := bworker.NewBWorker(10, bworker.WithRetry(2))
	defer bw.Shutdown()

	// Add a job to the worker pool.
	bw.Do(func() error {
		fmt.Println("Success without retry")
		return nil
	})

	// Add another job to the worker pool with error.
	retried := 0
	bw.Do(func() error {
		retried++

		// Return success on the last attempt
		if retried == 3 { // 1 (base) + 2 (retry)
			fmt.Println("Success with retry:", retried)
			return nil
		}

		fmt.Println("An error occur:", retried)
		return errors.New("an error")
	})

	// Wait for all jobs to finish executing.
	bw.Wait()
}
```

```text
Success without retry
An error occur: 1
An error occur: 2
Success with retry: 3
```

To use the options that BWorker provides, you can pass them to the `NewBWorker()` function. For example, the following
code creates a new worker pool with 4 workers and a job buffer of 10:

```go
bw := bworker.NewBWorker(4, bworker.WithJobBuffer(10))
```

The following code creates a new worker pool with 4 workers and a retry count of 3:

```go
bw := bworker.NewBWorker(4, bworker.WithRetry(3))
```

The following code creates a new worker pool with 4 workers and an error variable:

```go
var err error
bw := bworker.NewBWorker(4, bworker.WithError(&err))
```

The following code creates a new worker pool with 4 workers and a slice of error variables:

```go
var errs []error
bw := bworker.NewBWorker(4, bworker.WithErrors(&errs))
```

Once you have created a worker pool with options, you can use it just like a regular worker pool. The following code
shows how to submit a task to a worker pool with options:

```go
bw.Do(func () error {
	// ...
})
```

The following code shows how to wait for all of the tasks in a worker pool with options to be completed:

```go
bw.Wait()
```

The following code shows how to shutdown a worker pool with options:

```go
bw.Shutdown()
```

You can also use options to combine different features. For example, the following code creates a new worker pool with 4
workers, a job buffer of 10, and a retry count of 3:

```go
bw := bworker.NewBWorker(4, bworker.WithJobBuffer(10), bworker.WithRetry(3))
```

## Options

BWorker provides a number of options that can be used to customize the behavior of the worker pool. These options can be
passed to the `NewBWorker()` function.

The following options are currently available:

- **WithJobBuffer:** This option specifies the size of the job buffer. The job buffer is a queue of tasks that are
  waiting to be executed by a worker.
- **WithRetry:** This option specifies the number of times that a task will be retried if it fails.
- **WithError:** This option specifies a pointer to an error variable. If any of the tasks in the worker pool fail, the
  error will be stored in this variable.
- **WithErrors:** This option specifies a pointer to a slice of error variables. If any of the tasks in the worker pool
  fail, the errors will be stored in this slice.

## TODO

- Add BWorkerFlex documentation
- Add backoff retry
- Add fixed delay retry

## License

This project is licensed under the MIT License - see the LICENSE file for details.
