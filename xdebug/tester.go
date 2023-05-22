package xdebug

import "fmt"

type CaseTester map[int]int

func (c CaseTester) Run() bool {
	c[-1]++
	ok := c[c[-1]] > 0 || c[0] > 0
	if ok {
		fmt.Printf("\n\n\n>>Case %v is starting\n", c[-1])
	} else {
		fmt.Printf("\n\n\n>>Case %v is skipped\n", c[-1])
	}
	return ok
}
