package xdebug

import (
	"fmt"
	"strings"
)

type CaseTester map[int]int

func (c CaseTester) Run(message ...string) bool {
	c[-1]++
	ok := c[c[-1]] > 0 || c[0] > 0
	m := ""
	if len(message) > 0 {
		m = fmt.Sprintf("(%v)", strings.Join(message, " "))
	}
	if ok {
		fmt.Printf("\n\n\n>>Case%v %v is starting\n", m, c[-1])
	} else {
		fmt.Printf("\n\n\n>>Case%v %v is skipped\n", m, c[-1])
	}
	return ok
}
