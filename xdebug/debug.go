package xdebug

import "runtime"

func CallStack() string {
	buf := make([]byte, 102400)
	blen := runtime.Stack(buf, false)
	return string(buf[0:blen])
}
