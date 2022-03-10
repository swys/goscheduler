package goscheduler_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/swys/goscheduler"
)

type Task struct {
	count int64
}

func (t *Task) task() {
	atomic.StoreInt64(&t.count, atomic.LoadInt64(&t.count)+1)
	fmt.Printf("hello, I am a task\n")
}

func NewTask() *Task {
	return &Task{}
}

func TestGoScheduler(t *testing.T) {
	cases := []struct {
		description string
		interval    time.Duration
		Task        *Task
		sleep       time.Duration
		expected    int64
	}{
		{"task should run 8 times", time.Duration(1) * time.Second, NewTask(), time.Duration(8 * time.Second), 8},
		{"task should run 2 times", time.Duration(1) * time.Second, NewTask(), time.Duration(2 * time.Second), 2},
		{"task should run 0 times", time.Duration(10) * time.Second, NewTask(), time.Duration(1 * time.Second), 0},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t2 *testing.T) {
			scheduler := goscheduler.NewScheduler(tc.interval, tc.Task.task)
			go scheduler.SchedulerFunc()
			fmt.Printf("Sleeping for [%v] seconds...\n", tc.sleep)
			time.Sleep(tc.sleep)
			fmt.Printf("Stopping tasks...\n")
			scheduler.Done()
			if atomic.LoadInt64(&tc.Task.count) != tc.expected {
				t2.Errorf("Expected [%v] but got [%v]", tc.expected, atomic.LoadInt64(&tc.Task.count))
			}
		})
	}
}

func TestGoSchedulerDoneSafety(t *testing.T) {
	cases := []struct {
		description string
		scheduler   *goscheduler.Scheduler
		expected    bool
		callDone    int
	}{
		{"scheduler is initialized already but scheduler not running...calling Done() has no effect", goscheduler.NewScheduler(time.Duration(5)*time.Second, NewTask().task), true, 1},
		{"scheduler is NOT initialized yet...calling Done() multiple times has no effect", &goscheduler.Scheduler{}, false, 2},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t2 *testing.T) {
			for i := 0; i <= tc.callDone; i++ {
				tc.scheduler.Done()
			}
		})
	}
}

func TestIsRunning(t *testing.T) {
	cases := []struct {
		description string
		scheduler   *goscheduler.Scheduler
		expected    int64
	}{
		{"scheduler is running", goscheduler.NewScheduler(time.Duration(5)*time.Second, NewTask().task), 1},
		{"scheduler is NOT running", &goscheduler.Scheduler{}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t2 *testing.T) {
			defer tc.scheduler.Done()
			if tc.scheduler.SchedulerFunc != nil {
				go tc.scheduler.SchedulerFunc()
			}
			if tc.expected != tc.scheduler.IsRunning() {
				t2.Errorf("expected [%v] but got [%v]", tc.expected, tc.scheduler.IsRunning())
			}
		})
	}
}
