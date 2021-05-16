package xio

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestTimeoutListener(t *testing.T) {
	var err error
	var conn net.Conn
	{ //test timeout
		ln, _ := net.Listen("tcp", ":0")
		listener := NewTimeoutListener(ln, 100*time.Millisecond)
		listener.Delay = time.Millisecond
		var doneCount = 0
		go func() {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			conn.Read(make([]byte, 1024))
			doneCount++
		}()
		go func() {
			conn, err = net.Dial("tcp", ln.Addr().String())
			if err != nil {
				panic(err)
			}
			conn.Read(make([]byte, 1024))
			conn.Close()
			doneCount++
		}()
		time.Sleep(150 * time.Millisecond)
		if doneCount != 2 {
			t.Error("error")
			return
		}
		listener.Close()
	}
	{ //test read write close
		ln, _ := net.Listen("tcp", ":0")
		listener := NewTimeoutListener(ln, 100*time.Millisecond)
		listener.Delay = time.Millisecond
		var doneCount = 0
		go func() {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			io.Copy(conn, conn)
			conn.Close()
			doneCount++
		}()
		go func() {
			conn, err = net.Dial("tcp", ln.Addr().String())
			if err != nil {
				panic(err)
			}
			for {
				_, err = fmt.Fprintf(conn, "abc")
				if err != nil {
					break
				}
				buf := make([]byte, 3)
				_, err = conn.Read(buf)
				if err != nil {
					break
				}
				if string(buf) != "abc" {
					t.Error("error")
					return
				}
				time.Sleep(80 * time.Millisecond)
			}
			conn.Close()
			doneCount++
		}()
		time.Sleep(300 * time.Millisecond)
		if doneCount != 0 {
			t.Error("error")
			return
		}
		conn.Close()
		listener.Close()
	}
}
