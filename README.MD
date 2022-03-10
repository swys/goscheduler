goscheduler
===========

go program that takes a "task" function and a `time.Duration` interval and will call the task continuously until the program exists or `Done()` method is called on the scheduler object.

example
=======

```
package main

import (
  "fmt"
  "time"
  "github.com/swys/goscheduler"
)

func task() {
  fmt.Printf("I am a task\n")
}

func main() {
  // setup interval
  interval := time.Duration(1) * time.Second
  // get scheduler object and schedulerFunc (go routine used to run the task based on interval)
  scheduler, schedulerFunc := goscheduler.NewScheduler(interval, task)
  // start the interval to run the tasks
  go schedulerFunc()
  // sleep for 2 secs to simulate a later time when you want to stop the scheduler
  time.Sleep(2 * time.Second)
  // stop the scheduler
  scheduler.Done()
}
```

functions
=========

#### NewScheduler(interval time.Duration, task func()) (*Scheduler, func())

this creates a new scheduler object and a schedulerFunc which is a go routine used to start the scheduler. You must start the scheduler by running the returned schedulerFunc as a go routine like below :

```
scheduler, schedulerFunc := NewScheduler(interval, task)
go schedulerFunc()
```

##### arguments:
|argument  |type   |example   |description   |
|---|---|---|---|
|interval   |time.Duration   |time.Duration(1)*time.Second   |how often should your task run? |
|task   |func()   |func task() {...}   |function to run |

##### return:
|return  |type   |description |
|---|---|---|
|scheduler   |*Scheduler |scheduler object with `Done()` method for stopping scheduler|
|schedulerFunc   |func() |go routine used to start the scheduler|

methods
=======

#### Done()

call this method on the `Scheduler` object when you want to stop the scheduler from running. This method is safe to call as many time as you want, weather the Scheduler go routine is running or not or even if the scheduler has not been initialized yet with an interval and a task function.

```
scheduler, schedulerFunc := NewScheduler(interval, task)
go schedulerFunc()
// some time in the future when you want to stop the scheduler
scheduler.Done()
// from this point on, the task will not be called anymore
```

#### IsRunning() int64 

call this method when you want to find out if the Scheduler's go routine is running or not. This will return an `int64` and a value of `0` means the scheduler is NOT running. A value of `1` means the scheduler is running.

```
scheduler, schedulerFunc := NewScheduler(interval, task)
go schedulerFunc()
// find out if scheduler is running
isRunning := scheduler.IsRunning() // returns 1
```

test
====

```
make test
```

license
=======

MIT
