// Package scheduler provides a simple, humans-friendly way to schedule the execution of the go function.
// It includes delay execution and periodic execution.
//
// Copyright (c) 2019, prprprus All rights reserved.
// Use of this source code is governed by a BSD-style .
// license that can be found in the LICENSE file.
package scheduler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"
)

const (
	// Delay represents job type, the job will be delayed execute once according to job sched
	Delay = "Delay"

	// Every represents job type, the job will be cycled execute according to job sched
	Every = "Every"

	// Key of job sched
	Second  = "Second"
	Minute  = "Minute"
	Hour    = "Hour"
	Day     = "Day"
	Weekday = "Weekday"
	Month   = "Month"

	// EveryRune is value of job sched, like "*" in cron, represents every
	// second/minute/hour/day/weekday/month.
	EveryRune = -1
)

var (
	// defaultJobSetSize default size for job set
	defaultJobSetSize = 5000

	// maxJobSetSize maximum size for job set
	maxJobSetSize = 10000

	// ErrPendingJob is returned when the pending job not exist
	ErrPendingJob = errors.New("pending job not exist")

	// ErrOverlength is returned when the job size over maxJobSetSize variable
	ErrOverlength = errors.New("job set size overlength")

	// ErrJobType is returned when the job type not exist,
	// job type is one of Delay and Every.
	ErrJobType = errors.New("job type not exist")

	// ErrJobSched is returned when the job sched not exist,
	// under normal circumstances, this error will not occur,
	// unless the key definition of job sched is incorrectly modified.
	ErrJobSched = errors.New("job sched not exist")

	// ErrTimeNegative is returned when the time argument is negative
	ErrTimeNegative = errors.New("time argument can not be negative")

	// ErrDupJobID is returned when the generateID generates the same id
	ErrDupJobID = errors.New("Duplicate job id")

	// ErrAlreadyComplayed is returned when cancel a completed job
	ErrAlreadyComplayed = errors.New("Job hash been completed")

	// ErrCancelJob is returned when time.Timer.Stop function occur error
	ErrCancelJob = errors.New("cancel job failed")

	// ErrRangeSecond is returned when Second method argument is not int
	ErrRangeSecond = errors.New("argument 0 <= n <= 59 in Second method")

	// ErrRangeMinute is returned when Minute method argument is not int
	ErrRangeMinute = errors.New("argument 0 <= n <= 59 in Minute method")

	// ErrRangeHour is returned when Hour method argument is not int
	ErrRangeHour = errors.New("argument 0 <= n <= 23 in Hour method")

	// ErrRangeDay is returned when Day method argument is not int
	ErrRangeDay = errors.New("argument 1 <= n <= 31 in Day method")

	// ErrRangeWeekday is returned when Weekday method argument is not int
	ErrRangeWeekday = errors.New("argument 0 <= n <= 6 in Weekday method")

	// ErrRangeMonth is returned when Month method argument is not int
	ErrRangeMonth = errors.New("argument 1 <= n <= 12 in Month method")

	// EmptyJobType represents an empty job type
	EmptyJobType = ""

	// EmptySched represents an empty job sched
	EmptySched = map[string]int{}

	// jobSet is a instance of JobSet
	jobSet = &JobSet{
		lock:         new(sync.Mutex),
		pendingSet:   map[string]*Job{},
		completedSet: map[string]bool{},
	}
)

// JobSet

// JobSet stores pending jobs and completed jobs and it is concurrent safly.
type JobSet struct {
	lock         *sync.Mutex     // ensure concurrent safe
	pendingSet   map[string]*Job // storage pending jobs
	completedSet map[string]bool // storage completed jobs
}

// setJobDone When the job function is executed then set job done.
func (js *JobSet) setJobDone(id string) {
	js.lock.Lock()
	defer js.lock.Unlock()

	job := js.pendingSet[id]

	// note: ignore with job type is Every
	if job.Type != Every {
		if _, ok := js.completedSet[id]; ok {
			panic(ErrDupJobID)
		}

		delete(js.pendingSet, id)
		js.completedSet[id] = true
	}
}

// Job

// JobTimer is the wrapper for time.Timer, one job corresponds to a JobTimer.
type JobTimer struct {
	ID    string      // unique id
	timer *time.Timer // wrapper time.Timer
}

// JobTicker is the wrapper for time.Ticker, one job corresponds to a JobTicker.
type JobTicker struct {
	ID     string       // unique id
	ticker *time.Ticker // wrapper time.Ticker
}

// Job is an abstraction of a scheduling task.
type Job struct {
	ID   string // unique id
	Type string // job type

	// Sched is a job sched, like cron style but the order of time is not
	// fixed, can be arranged and combined at will.
	Sched map[string]int

	fn      interface{}   // job function
	args    []interface{} // function args
	JTimer  *JobTimer     // JobTimer
	JTicker *JobTicker    // JobTicker
}

// Second method set Second key for job sched.
func (j *Job) Second(n int) *Job {
	if n < 0 || n > 59 {
		panic(ErrRangeSecond)
	}

	j.Sched[Second] = n
	return j
}

// Minute method set Minute key for job sched.
func (j *Job) Minute(n int) *Job {
	if n < 0 || n > 59 {
		panic(ErrRangeMinute)
	}

	j.Sched[Minute] = n
	return j
}

// Hour method set Hour key for job sched.
func (j *Job) Hour(n int) *Job {
	if n < 0 || n > 23 {
		panic(ErrRangeHour)
	}

	j.Sched[Hour] = n
	return j
}

// Day method set Day key for job sched.
func (j *Job) Day(n int) *Job {
	if n < 1 || n > 31 {
		panic(ErrRangeDay)
	}

	j.Sched[Day] = n
	return j
}

// Weekday method set Weekday key for job sched.
func (j *Job) Weekday(n int) *Job {
	if n < 0 || n > 6 {
		panic(ErrRangeWeekday)
	}

	j.Sched[Weekday] = n
	return j
}

// Month method set Month key for job sched.
func (j *Job) Month(n int) *Job {
	if n < 1 || n > 12 {
		panic(ErrRangeMonth)
	}

	j.Sched[Month] = n
	return j
}

// Do according to the job type and job sched execute job.
func (j *Job) Do(fn interface{}, args ...interface{}) (jobID string) {
	j.fn = fn
	j.args = args

	switch j.Type {
	case Delay:
		// convert to second. Not support Weekday and Month
		var second int
		for k := range j.Sched {
			switch k {
			case Second:
				second = j.Sched[Second]
			case Minute:
				second = j.Sched[Minute] * 60
			case Hour:
				second = j.Sched[Hour] * 60 * 60
			case Day:
				second = j.Sched[Day] * 60 * 60 * 24
			default:
				panic(ErrJobSched)
			}
		}
		// initial job.JTimer (note: can not put it in a new goroutine)
		j.JTimer = new(JobTimer)
		j.JTimer.ID = generateID()
		j.JTimer.timer = time.NewTimer(time.Duration(second) * time.Second)
		go func() {
			// wait...
			<-j.JTimer.timer.C
			// run job function
			j.run()
			// set job done
			jobSet.setJobDone(j.ID)
		}()
	case Every:
		// initial job.JTicker (note: also can not put it in a new goroutine)
		j.JTicker = new(JobTicker)
		j.JTicker.ID = generateID()
		j.JTicker.ticker = time.NewTicker(1 * time.Second)
		go func() {
			// begin ticktock...
			for t := range j.JTicker.ticker.C {
				_ = t
				if (j.Sched[Second] == -1 || j.Sched[Second] == time.Now().Second()) &&
					(j.Sched[Minute] == -1 || j.Sched[Minute] == time.Now().Minute()) &&
					(j.Sched[Hour] == -1 || j.Sched[Hour] == time.Now().Hour()) &&
					(j.Sched[Day] == -1 || j.Sched[Day] == time.Now().Day()) &&
					(j.Sched[Weekday] == -1 || j.Sched[Weekday] == int(time.Now().Weekday())) &&
					(j.Sched[Month] == -1 || j.Sched[Month] == int(time.Now().Month())) {
					// run job function
					j.run()
					// set job done
					jobSet.setJobDone(j.ID)
				}
			}
		}()
	default:
		panic(ErrJobType)
	}

	return j.ID
}

// run funtion of job by reflect.
func (j *Job) run() {
	rFn := reflect.ValueOf(j.fn)
	rArgs := make([]reflect.Value, len(j.args))
	for i, v := range j.args {
		rArgs[i] = reflect.ValueOf(v)
	}

	// retry
	defer func() {
		if err := recover(); err != nil {
			time.Sleep(5 * time.Second) // wait for five seconds for now
			rFn.Call(rArgs)
		}
	}()

	rFn.Call(rArgs)
}

// Scheduler

// Scheduler is responsible for scheduling jobs.
type Scheduler struct {
	jobSetSize int     // custom size for job set, can not overlength maxJobSetSize
	js         *JobSet // JobSet
}

// NewScheduler new Scheduler instance.
func NewScheduler(jss int) (*Scheduler, error) {
	if jss > maxJobSetSize {
		return nil, ErrOverlength
	}
	if jss <= 0 {
		jss = defaultJobSetSize
	}

	s := &Scheduler{
		jobSetSize: jss,
		js:         jobSet,
	}
	return s, nil
}

// Delay method schedule job with Delay mode.
func (s *Scheduler) Delay() *Job {
	s.js.lock.Lock()
	defer s.js.lock.Unlock()

	// temporarily handle like this
	if len(s.js.pendingSet) >= 10000 {
		panic("pending set is full")
	}

	// create job
	id := generateID()
	j := &Job{
		ID:    id,
		Type:  Delay,
		Sched: InitJobSched(Delay),
	}

	// put in pending job set
	s.js.pendingSet[id] = j
	return j
}

// Every method schedule job with Every mode.
func (s *Scheduler) Every() *Job {
	s.js.lock.Lock()
	defer s.js.lock.Unlock()

	// temporarily handle like this
	if len(s.js.pendingSet) >= 10000 {
		panic("pending set is full")
	}

	// create job
	id := generateID()
	j := &Job{
		ID:   id,
		Type: Every,
		// Sched[...] = -1 <=> cron *
		Sched: InitJobSched(Every),
	}

	// put in pending job set
	s.js.pendingSet[id] = j
	return j
}

// PendingJob get pending job by id.
func (s *Scheduler) PendingJob(id string) (*Job, error) {
	s.js.lock.Lock()
	defer s.js.lock.Unlock()

	if job, ok := s.js.pendingSet[id]; ok {
		return job, nil
	}
	return nil, ErrPendingJob
}

// JobType get job type.
func (s *Scheduler) JobType(id string) (string, error) {
	job, err := s.PendingJob(id)
	if err != nil {
		return EmptyJobType, err
	}

	return job.Type, nil
}

// JobSched get job sched.
func (s *Scheduler) JobSched(id string) (map[string]int, error) {
	job, err := s.PendingJob(id)
	if err != nil {
		return EmptySched, err
	}

	return job.Sched, nil
}

// JobDone Check if the job is completed.
func (s *Scheduler) JobDone(id string) (bool, error) {
	s.js.lock.Lock()
	defer s.js.lock.Unlock()

	if _, ok := s.js.completedSet[id]; ok {
		return true, nil
	}
	return false, nil
}

// CancelJob can cancel the job before scheduling.
func (s *Scheduler) CancelJob(id string) error {
	// can not cancel a completed job
	if _, ok := s.js.completedSet[id]; ok {
		return ErrAlreadyComplayed
	}
	// can not cancel a nonexistent job
	if _, ok := s.js.pendingSet[id]; !ok {
		return ErrPendingJob
	}

	// cancel by job type
	job := s.js.pendingSet[id]
	switch job.Type {
	case Delay:
		ok := job.JTimer.timer.Stop()
		if ok {
			return nil
		}
		return ErrCancelJob
	case Every:
		job.JTicker.ticker.Stop()
		return nil
	default:
		return ErrJobType
	}
}

// util

// generateID generate job id
func generateID() string {
	h := md5.New()
	io.WriteString(h, time.Now().String())
	id := fmt.Sprintf("%x", h.Sum(nil))
	return id
}

// InitJobSched initiate job sched by job type
func InitJobSched(jobType string) map[string]int {
	if jobType == Delay {
		return map[string]int{}
	}
	return map[string]int{
		Second:  EveryRune,
		Minute:  EveryRune,
		Hour:    EveryRune,
		Day:     EveryRune,
		Weekday: EveryRune,
		Month:   EveryRune,
	}
}
