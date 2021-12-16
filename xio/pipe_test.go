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
		fmt.Printf("--->write0\n")
		fmt.Fprintf(a, "123")
		fmt.Printf("--->write1\n")
		time.Sleep(10 * time.Millisecond)
		b.Close()
		fmt.Printf("--->close0\n")
	}()
	var n int
	buf := make([]byte, 1024)
	n, err = a.Read(buf)
	if err != nil || n != 3 || "abc" != string(buf[0:3]) {
		t.Error(err)
		return
	}
	fmt.Printf("--->read0\n")
	n, err = a.Read(buf)
	if err != nil || n != 3 || "123" != string(buf[0:3]) {
		t.Error(err)
		return
	}
	fmt.Printf("--->read1\n")
	_, err = a.Read(buf)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Printf("--->close0\n")
	a.Close()
	a.Read(nil)
	a.Write(nil)
}

func TestPipe2(t *testing.T) {
	for i := 0; i < 10000; i++ {
		a, _, err := Pipe()
		if err != nil {
			t.Error(err)
			return
		}
		a.Close()
	}
}

func TestPipedConne(t *testing.T) {
	a, b, err := CreatePipedConn()
	if err != nil {
		t.Error(err)
		return
	}
	a.RemoteAddr()
	a.LocalAddr()
	a.SetDeadline(time.Now())
	a.SetReadDeadline(time.Now())
	a.SetWriteDeadline(time.Now())
	a.Network()
	fmt.Printf("-->%v\n", a)
	b.Close()
}

func TestPipedListener(t *testing.T) {
	listener := NewPipedListener()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			go io.Copy(conn, conn)
		}
	}()
	conn, err := listener.Dial()
	if err != nil {
		return
	}
	fmt.Fprintf(conn, "abc")
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	if string(buf[0:n]) != "abc" {
		t.Error("error")
		return
	}
	conn.Close()
	listener.Close()
	time.Sleep(100 * time.Millisecond)
	listener.Addr()
	listener.Network()
	fmt.Printf("listener %v\n", listener)
}
