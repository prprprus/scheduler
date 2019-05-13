# scheduler

![build status](https://travis-ci.org/prprprus/scheduler.svg?branch=master)
[![codecov](https://codecov.io/gh/prprprus/scheduler/branch/master/graph/badge.svg)](https://codecov.io/gh/prprprus/scheduler)
[![godoc](https://img.shields.io/badge/godoc-godoc-blue.svg)](https://godoc.org/github.com/prprprus/scheduler)
[![license](https://img.shields.io/badge/license-license-yellow.svg)](https://github.com/prprprus/scheduler/blob/master/LICENSE)

[英文文档](https://github.com/prprprus/scheduler)

## 介绍

scheduler 是 Go 语言实现的作业调度工具包。它提供了一种简单、人性化的方式去调度 Go 函数，包括延迟和周期性两种调度方式。

灵感来源于 Linux [cron](https://opensource.com/article/17/11/how-use-cron-linux) 和 Python [schedule](https://github.com/dbader/schedule)。

## 功能

- 延迟执行，精确到一秒钟
- 周期性执行，精确到一秒钟，类似 cron 的风格，但是更加的灵活
- 取消 job
- 失败重试（暂时重试一次）

## 安装

```
go get github.com/prprprus/scheduler
```

## 例子

job 函数

```Go
func task1(name string, age int) {
	fmt.Printf("run task1, with arguments: %s, %d\n", name, age)
}

func task2() {
	fmt.Println("run task2, without arguments")
}
```

### 延迟调度

延迟调度支持四种模式：按秒、分、小时、天。

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
	jobID := s.Delay().Second(1).Do(task1, "prprprus", 23)

	// delay with 1 minute, job function without arguments
	jobID = s.Delay().Minute(1).Do(task2)

	// delay with 1 hour
	jobID = s.Delay().Hour(1).Do(task2)

	// Note: execute immediately
	jobID = s.Delay().Do(task2)

	// cancel job
	jobID = s.Delay().Day(1).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("cancel delay job success")
	}
}
```

### 周期性调度

类似 cron 的风格，同样会包括秒、分、小时、天、星期、月，但是它们之间的顺序和数量不需要固定成一个死格式。你可以按照你的个人喜好去进行排列组合。例如，`Second(3).Minute(35).Day(6)` 和 `Minute(35).Day(6).Second(3)` 的效果是一样的。不需要再去记格式了！🎉👏

但是为了可读性，推荐按照从小到大（或者从大到小）的顺序使用。

注意：`Day()` 和 `Weekday()` 避免同时出现，除非你清楚知道这天是星期几。

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

	// Specifies time to execute periodically
	jobID = s.Every().Second(45).Minute(20).Hour(13).Day(23).Weekday(3).Month(6).Do(task1, "prprprus", 23)
	jobID = s.Every().Second(15).Minute(40).Hour(16).Weekday(4).Do(task2)
	jobID = s.Every().Second(1).Do(task1, "prprprus", 23)

	// Note: execute immediately
	jobID = s.Every().Do(task2)

	jobID = s.Every().Second(1).Minute(1).Hour(1).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("cancel periodically job success")
	}
}
```

## Documentation

[完整的文档](https://godoc.org/github.com/prprprus/scheduler)

## Contribution

非常感谢你对该项目感兴趣，你的帮助对我来说是非常宝贵的。你可以提交 issue、pull requests 以及 fork，建议在 pull requests 之前先提交一个 issue 哈。

## License

[LICENSE](https://github.com/prprprus/scheduler/blob/master/LICENSE) 详情.

✨彩蛋✨：该项目是 Github 上第一个名字叫做 scheduler 的 Go 项目。👻