package task

import (
	"fmt"
	"sync"

	"github.com/codingeasygo/util/debug"
)

var Shared = NewRunner()

func Call(state interface{}, call func(state interface{})) {
	Shared.Call(state, call)
}

type Task struct {
	State  interface{}
	Caller func(state interface{})
}

type Runner struct {
	max     int64
	running int64
	queue   chan *Task
	waiter  sync.WaitGroup
	locker  sync.RWMutex
}

func NewRunner() (runner *Runner) {
	runner = &Runner{
		queue:  make(chan *Task, 1024),
		waiter: sync.WaitGroup{},
		locker: sync.RWMutex{},
		max:    3,
	}
	return
}

func (r *Runner) Max(max int) {
	r.max = int64(max)
}

func (r *Runner) Stop() {
	r.locker.Lock()
	for i := int64(0); i < r.running; i++ {
		r.queue <- nil
	}
	r.running = 0
	r.locker.Unlock()
	r.waiter.Wait()
}

func (r *Runner) run() {
	defer r.waiter.Done()
	running := true
	for running {
		t := <-r.queue
		if t == nil {
			running = false
			break
		}
		func(t_ *Task) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("Runner call panic with %v, callstack is \n%v", err, debug.CallStatck())
				}
			}()
			t_.Caller(t_.State)
		}(t)
		r.locker.Lock()
		if r.running > r.max {
			r.running--
			running = false
		}
		r.locker.Unlock()
	}
}

func (r *Runner) Call(state interface{}, call func(state interface{})) {
	r.locker.Lock()
	if r.running < r.max {
		r.running++
		r.waiter.Add(1)
		go r.run()
	}
	r.locker.Unlock()
	r.queue <- &Task{State: state, Caller: call}
}
