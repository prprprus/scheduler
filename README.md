![scheduler6.png](https://i.loli.net/2019/09/21/CbpFx7TIvSNM1EP.png)

![build status](https://travis-ci.org/prprprus/scheduler.svg?branch=master)
[![codecov](https://codecov.io/gh/prprprus/scheduler/branch/master/graph/badge.svg)](https://codecov.io/gh/prprprus/scheduler)
[![godoc](https://img.shields.io/badge/godoc-godoc-blue.svg)](https://godoc.org/github.com/prprprus/scheduler)
[![license](https://img.shields.io/badge/license-license-yellow.svg)](https://github.com/prprprus/scheduler/blob/master/LICENSE)

[中文文档](https://github.com/prprprus/scheduler/blob/master/README-zh.md)

## Introduction

The scheduler is a job scheduling package for Go. It provides a simple, human-friendly way to schedule the execution of the go function and includes delay and periodic.

Inspired by Linux [cron](https://opensource.com/article/17/11/how-use-cron-linux) and Python [schedule](https://github.com/dbader/schedule).

## Features

- Delay execution, accurate to a second
- Periodic execution, accurate to a second, like the cron style but more flexible
- Cancel job
- Failure retry

## Installation

```
go get github.com/prprprus/scheduler
```

## Example

job function

```Go
func task1(name string, age int) {
	fmt.Printf("run task1, with arguments: %s, %d\n", name, age)
}

func task2() {
	fmt.Println("run task2, without arguments")
}
```

### Delay

Delayed supports four modes: seconds, minutes, hours, and days.

As a special case, the task will be executed immediately via `s.Delay().Do(task)` .

```Go
package main

import (
    "fmt"

    "github.com/prprprus/scheduler"
)

func main() {
	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		panic(err) // just example
	}

	// delay with 1 second, job function with arguments
	s.Delay().Second(1).Do(task1, "prprprus", 23)

	// delay with 1 minute, job function without arguments
	s.Delay().Minute(1).Do(task2)

	// delay with 1 hour
	s.Delay().Hour(1).Do(task2)

	// special: execute immediately
	s.Delay().Do(task2)

	// cancel job
	jobID := s.Delay().Day(1).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("cancel delay job success")
	}
}
```

### Every

Like the cron style, it also includes seconds, minutes, hours, days, weekday, and months, but the order and number are not fixed. You can freely arrange and combine them according to your own preferences. For example, the effects of `Second(3).Minute(35).Day(6)` and `Minute(35).Day(6).Second(3)` are the same. No need to remember the format! 🎉👏

But for the readability, recommend the chronological order from small to large (or large to small).

Note: `Day()` and `Weekday()` avoid simultaneous occurrences unless you know that the day is the day of the week.

As a special case, the task will be executed once per second via `s.Every().Do(task)`.

```Go
package main

import (
    "fmt"

    "github.com/prprprus/scheduler"
)

func main() {
	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		panic(err)
	}

	// Specifies time to execute periodically
	s.Every().Second(45).Minute(20).Hour(13).Day(23).Weekday(3).Month(6).Do(task1, "prprprus", 23)
	s.Every().Second(15).Minute(40).Hour(16).Weekday(4).Do(task2)
	s.Every().Second(1).Do(task1, "prprprus", 23)

	// special: executed once per second
	s.Every().Do(task2)

	// cancel job
	jobID := s.Every().Second(1).Minute(1).Hour(1).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("cancel periodically job success")
	}
}
```

## Documentation

[Full documentation](https://godoc.org/github.com/prprprus/scheduler)

## Contribution

Thank you for your interest in the contribution of the scheduler, your help and contribution is very valuable.

You can submit an issue and pull requests and fork, please submit an issue before submitting pull requests.

## License

See [LICENSE](https://github.com/prprprus/scheduler/blob/master/LICENSE) for more information.
