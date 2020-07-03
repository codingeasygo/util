package xtest

import (
	"fmt"
	"testing"
	"time"
)

func TestPerf(t *testing.T) {
	DoPerfV(20, 10, "", func(v int) {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("doing->%d\n", v)
	})
	used, err := DoPerf(1000, "t.log", func(v int) {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("doing->%d\n", v)
		fmt.Println(v)
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("used->", used)
}

// func TestAutoPerf(t *testing.T) {
// 	perf := NewPerf()
// 	perf.ShowState = true
// 	used, max, avg, err := perf.AutoExec(1000, 10, 10, "", 100,
// 		func(idx int, state Perf) error {
// 			fmt.Printf("running->%d->%v\n", idx, converter.JSON(state))
// 			if state.Running < 100 {
// 				return nil
// 			}
// 			fmt.Println("--->full")
// 			return FullError
// 		}, func(v int) error {
// 			time.Sleep(2000 * time.Millisecond)
// 			fmt.Printf("doing->%d\n", v)
// 			return nil
// 		})
// 	fmt.Println("used->", used, max, avg, err)
// }
