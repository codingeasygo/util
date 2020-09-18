package xtime

import "time"

func Now() int64 {
	return time.Now().Local().UnixNano() / 1e6
}
