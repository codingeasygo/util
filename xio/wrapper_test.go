package xio

import (
	"fmt"
	"testing"
	"time"
)

func TestConnWrapper(t *testing.T) {
	conn := NewConnWrapper(NewCombinedReadWriteCloser(nil, nil, nil))
	conn.LocalAddr()
	conn.RemoteAddr()
	conn.SetDeadline(time.Now())
	conn.SetReadDeadline(time.Now())
	conn.SetWriteDeadline(time.Now())
	conn.Network()
	fmt.Println(conn.String())
}
