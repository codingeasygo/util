package task

import (
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	Shared.Max(3)
	Call(1, func(state interface{}) {
		fmt.Printf("--->%v\n", state)
	})
	Shared.Stop()
}
