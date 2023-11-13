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

## Worker Types

### 1. [Pool Worker](https://github.com/bearaujus/bworker/blob/master/pool)

An BWorker instance with **specified** concurrency level. Features:

- **WithJobPoolSize**: Set the size of the job pool size.
- **WithWorkerStartupDelay**: Set the delay interval between the startup of each worker in the worker pool.
- **WithRetry**: Set the number of times to retry a failed job.
- **WithError**: Set a pointer to an error variable that will be populated if any job fails.
- **WithErrors**: Set a pointer to a slice of error variables that will be populated if any job fails.

### 2. [Flex Worker](https://github.com/bearaujus/bworker/blob/master/flex)

An BWorker instance with **unlimited** concurrency level. Features:

- **WithRetry**: Set the number of times to retry a failed job.
- **WithError**: Set a pointer to an error variable that will be populated if any job fails.
- **WithErrors**: Set a pointer to a slice of error variables that will be populated if any job fails.

## TODO

- Add `PoolWorker` documentation & example
- Add `FlexWorker` documentation & example
- Add internal test cases
- Support backoff retry
- Support fixed delay retry
- Function `Do` and `DoSimple` at `PoolWorker` sometimes will block the thread,
  so we need to handle this by adding `DoWithContext` & `DoSimpleWithContext`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
