// Copyright (c) 2019, prprprus All rights reserved.
// Use of this source code is governed by a BSD-style .
// license that can be found in the LICENSE file.

package scheduler

import (
	"reflect"
	"testing"
	"time"
)

func task1(name, age string, res *[]string) {
	*res = append(*res, name, age)
}

func task2() {
	s := "hello-world, task2, without arguments"
	_ = s
}

// JobSet

func TestJobDone(t *testing.T) {
	s, _ := NewScheduler(10)
	jobID := s.Delay().Second(0).Do(task2)
	time.Sleep(1 * time.Second)
	ok, _ := s.JobDone(jobID)
	if !ok {
		t.Errorf("job has been completed, JobDone method error")
	}

	jobID = s.Delay().Hour(1).Do(task2)
	time.Sleep(1 * time.Second)
	ok, _ = s.JobDone(jobID)
	if ok {
		t.Errorf("job not completed, JobDone method error")
	}
}

// Scheduler

func TestNewScheduler(t *testing.T) {
	_, err := NewScheduler(-1)
	if err != nil {
		t.Errorf("maxJobSetSize can be negative")
	}
	_, err = NewScheduler(10001)
	if err == nil {
		t.Errorf("maxJobSetSize overlength")
	}
}

func TestPendingJob(t *testing.T) {
	s, _ := NewScheduler(10)
	jobID := s.Delay().Minute(10).Do(task2)
	_, err := s.PendingJob(jobID)
	if err != nil {
		t.Errorf("pending job should be exists")
	}

	_, err = s.PendingJob(generateID())
	if err == nil {
		t.Errorf("pending job should not exists")
	}
}

func TestJobType(t *testing.T) {
	s, _ := NewScheduler(10)
	jobID := s.Delay().Day(3).Do(task2)
	_, err := s.JobType(jobID)
	if err != nil {
		t.Errorf("job type should be exists")
	}

	_, err = s.JobType(generateID())
	if err == nil {
		t.Errorf("job type should not exists")
	}
}

func TestSched(t *testing.T) {
	s, _ := NewScheduler(10)
	jobID := s.Delay().Day(1).Do(task2)
	_, err := s.JobSched(jobID)
	if err != nil {
		t.Errorf("job sched should be exists")
	}

	_, err = s.JobSched(generateID())
	if err == nil {
		t.Errorf("job sched should not exists")
	}
}

func TestCancelJob(t *testing.T) {
	s, _ := NewScheduler(10)

	// cancel completed job
	jobID := s.Delay().Second(0).Do(task2)
	time.Sleep(1 * time.Second)
	err := s.CancelJob(jobID)
	if err == nil {
		t.Errorf("job completed, can not cancel")
	}

	// cancel nonexistent job
	err = s.CancelJob(generateID())
	if err == nil {
		t.Errorf("job not exists, can not cancel")
	}

	jobID = s.Delay().Minute(30).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		t.Errorf("job should be cancel")
	}
}

// Job

func TestSecond(t *testing.T) {
	defer func() {
		if err := recover(); err != nil && err == ErrRangeSecond {
			return
		}
	}()

	s, _ := NewScheduler(10)
	j := s.Every().Second(233)
	if j.Sched[Second] != 233 {
		t.Errorf("set second error")
	}

	// panic
	j.Second(-1)
}

func TestMinute(t *testing.T) {
	defer func() {
		if err := recover(); err != nil && err == ErrRangeMinute {
			return
		}
	}()

	s, _ := NewScheduler(10)
	j := s.Every().Minute(59)
	if j.Sched[Minute] != 59 {
		t.Errorf("set minute error")
	}

	// panic
	j.Minute(-1)
}

func TestHour(t *testing.T) {
	defer func() {
		if err := recover(); err != nil && err == ErrRangeHour {
			return
		}
	}()

	s, _ := NewScheduler(10)
	j := s.Every().Hour(12)
	if j.Sched[Hour] != 12 {
		t.Errorf("set hour error")
	}

	// panic
	j.Hour(-1)
}

func TestDay(t *testing.T) {
	defer func() {
		if err := recover(); err != nil && err == ErrRangeDay {
			return
		}
	}()

	s, _ := NewScheduler(10)
	j := s.Every().Day(24)
	if j.Sched[Day] != 24 {
		t.Errorf("set day error")
	}

	// panic
	j.Day(-1)
}

func TestWeekday(t *testing.T) {
	defer func() {
		if err := recover(); err != nil && err == ErrRangeWeekday {
			return
		}
	}()

	s, _ := NewScheduler(10)
	j := s.Every().Weekday(5)
	if j.Sched[Weekday] != 5 {
		t.Errorf("set weekday error")
	}

	// panic
	j.Weekday(-1)
}

func TestMonth(t *testing.T) {
	defer func() {
		if err := recover(); err != nil && err == ErrRangeMonth {
			return
		}
	}()

	s, _ := NewScheduler(10)
	j := s.Every().Month(1)
	if j.Sched[Month] != 1 {
		t.Errorf("set month error")
	}

	// panic
	j.Month(-1)
}

func TestDo(t *testing.T) {
	// Delay with arguments
	res1 := []string{}
	res2 := []string{"tiger", "23"}
	s, _ := NewScheduler(10)
	s.Delay().Second(1).Do(task1, "tiger", "23", &res1)
	time.Sleep(2 * time.Second)
	for i, v := range res1 {
		if v != res2[i] {
			t.Errorf("Do method error with Delay")
		}
	}

	// Every with arguments
	res1 = []string{}
	res2 = []string{"cat", "5"}
	s, _ = NewScheduler(10)
	jobID := s.Every().Second(1).Do(task1, "cat", "5", &res1)
	time.Sleep(2 * time.Second)
	err := s.CancelJob(jobID)
	if err != nil {
		panic(err)
	}
	for i, v := range res1 {
		if v != res2[i] {
			t.Errorf("Do method error with Every")
		}
	}
}

// util

func TestInitJobSched(t *testing.T) {
	// Delay
	s, _ := NewScheduler(10)
	j := s.Delay()
	if reflect.TypeOf(j.Sched) != reflect.TypeOf(map[string]int{}) && len(j.Sched) == 0 {
		t.Errorf("Initial Delay job sched failed")
	}

	// Every
	s, _ = NewScheduler(10)
	j = s.Every()
	for k, v := range j.Sched {
		if k != Second && k != Minute && k != Hour && k != Day && k != Weekday && k != Month {
			t.Errorf("Initial Every job sched failed")
		}
		if v != -1 {
			t.Errorf("Initial Every job sched failed")
		}
	}
}
