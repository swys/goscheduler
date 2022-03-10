package goscheduler

import (
	"sync/atomic"
	"time"
)

//
// Structs
//

// Scheduler struct
type Scheduler struct {
	done          chan bool
	SchedulerFunc func()
	running       int64
}

//
// Methods
//

// Done method : with close the scheduler
func (s *Scheduler) Done() {
	if s.IsRunning() == 1 {
		s.done <- true
		atomic.StoreInt64(&s.running, 0)
	}
}

// Run method informs Scheduler object that the SchedulerFunc started running
func (s *Scheduler) run() {
	atomic.StoreInt64(&s.running, 1)
}

// IsRunning method returns current state of Scheduler's running field, which indicates if the Scheduler's SchedulerFunc is currently running or not
func (s *Scheduler) IsRunning() int64 {
	// allow some delay before checking in case this check happen directly after the go routine is called
	// this delay should give enough time for the go routine to start
	time.Sleep(time.Duration(100 * time.Millisecond))
	return atomic.LoadInt64(&s.running)
}

// SchedulerFunc method is a function that takes done channel, ticker, and task and returns a func() to be called as a go routine to start the scheduler
func schedulerFunc(done chan bool, ticker *time.Ticker, task func(), s *Scheduler) func() {
	return func() {
		defer close(done)
		s.run()
		for {
			select {
			case <-ticker.C:
				task()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}
}

// NewScheduler function : run task at specified intervals until quit channel is closed
// Returns a quit channel and a go routine to run the interval
func NewScheduler(interval time.Duration, task func()) *Scheduler {
	done := make(chan bool)
	ticker := time.NewTicker(interval)
	scheduler := &Scheduler{
		done:    done,
		running: 0,
	}
	scheduler.SchedulerFunc = schedulerFunc(done, ticker, task, scheduler)
	return scheduler
}
