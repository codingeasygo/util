package xio

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func TestPipe(t *testing.T) {
	a, b, err := Pipe()
	if err != nil {
		t.Error(err)
		return
	}
	go io.Copy(b, b)
	go func() {
		fmt.Fprintf(a, "abc")
		time.Sleep(10 * time.Millisecond)
		a.Close()
	}()
	buf := make([]byte, 1024)
	n, err := a.Read(buf)
	if err != nil || n != 3 || "abc" != string(buf[0:3]) {
		t.Error(err)
		return
	}
	_, err = a.Read(buf)
	if err == nil {
		t.Error(err)
		return
	}
}
