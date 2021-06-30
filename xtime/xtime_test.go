package xtime

import (
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {
	fmt.Printf("Now-->%v\n", Now())
	fmt.Printf("TimeStartOfToday-->%v\n", TimeStartOfToday())
	fmt.Printf("TimeStartOfWeek-->%v\n", TimeStartOfWeek())
	fmt.Printf("TimeStartOfMonth-->%v\n", TimeStartOfMonth())
}
