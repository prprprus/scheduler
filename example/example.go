package main

import (
	"fmt"
	"time"

	"github.com/prprprus/scheduler"
)

func task1(name string, age int) {
	fmt.Printf("run task1, with arguments: %s, %d\n", name, age)
}

func task2() {
	fmt.Println("run task2, without arguments")
}

func main() {
	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		panic(err) // just example
	}

	fmt.Println("-------------------------------------------------")
	fmt.Println()

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

	time.Sleep(3 * time.Second)
	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

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

	fmt.Println()
	fmt.Println("--------------------------------------------------")
}
