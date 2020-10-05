package xio

import (
	"fmt"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	echo := NewEchoConn()
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
	//
	conna, connb, _ := CreatePipedConn()
	dialer := NewEchoDialer()
	piper, _ := dialer.DialPiper("xxx", 1024)
	waiter := make(chan int, 1)
	go func() {
		piper.PipeConn(connb, "xxx")
		piper.Close()
		waiter <- 1
	}()
	fmt.Fprintf(conna, "abc")
	err = FullBuffer(conna, buffer, 3, nil)
	if err != nil || string(buffer[0:3]) != "abc" {
		t.Error(err)
		return
	}
	conna.Close()
	<-waiter
	dialer.Dial("", "")
}
