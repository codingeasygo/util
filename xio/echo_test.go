package xio

import (
	"fmt"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	echo, _ := NewEchoConn()
	buffer := make([]byte, 1024)
	go func() {
		fmt.Fprintf(echo, "%v", "abc")
	}()
	n, err := echo.Read(buffer)
	if err != nil || n != 3 {
		t.Error(err)
		return
	}
	echo.Close()
	echo.LocalAddr()
	echo.RemoteAddr()
	echo.SetDeadline(time.Now())
	echo.SetReadDeadline(time.Now())
	echo.SetWriteDeadline(time.Now())
	echo.Network()
	fmt.Println(echo.String())
}
