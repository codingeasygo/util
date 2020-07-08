package debug

import "runtime"

func CallStatck() string {
	buf := make([]byte, 102400)
	blen := runtime.Stack(buf, false)
	return string(buf[0:blen])
}
