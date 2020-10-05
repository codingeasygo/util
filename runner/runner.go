package runner

import (
	"fmt"
	"log"
	"time"
)

//ErrNotTask is error for not task
var ErrNotTask = fmt.Errorf("not task")

//NamedRunner will run call by delay
func NamedRunner(name string, delay time.Duration, running *bool, call func() error) {
	log.Printf("%v is starting", name)
	var finishCount = 0
	for *running {
		err := call()
		if err == nil {
			finishCount++
			continue
		}
		if err != ErrNotTask {
			log.Printf("%v is fail with %v", name, err)
		} else if finishCount > 0 {
			log.Printf("%v is having %v finished", name, finishCount)
		}
		finishCount = 0
		time.Sleep(delay)
	}
	log.Printf("%v is stopped", name)
}
