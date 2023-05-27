package proxy

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"testing"

	"github.com/codingeasygo/util/xio"
)

func runEchoServer(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go io.Copy(conn, conn)
	}
}

func TestForward(t *testing.T) {
	go runEchoServer("127.0.0.1:13200")
	forward := NewForward("Test")
	{
		listenURL, _ := url.Parse("socks://127.0.0.1:0")
		_, err := forward.StartForward("socks", listenURL, "tcp://${HOST}")
		if err != nil {
			t.Error(err)
			return
		}
		forward.StopForward("socks")
	}
	{
		listenURL, _ := url.Parse("proxy://127.0.0.1:0")
		_, err := forward.StartForward("proxy", listenURL, "tcp://${HOST}")
		if err != nil {
			t.Error(err)
			return
		}
		forward.StopForward("proxy")
	}
	{
		listenURL, _ := url.Parse("tcp://127.0.0.1:13300")
		_, err := forward.StartForward("abc", listenURL, "tcp://127.0.0.1:13200")
		if err != nil {
			t.Error(err)
			return
		}
		conn, err := net.Dial("tcp", "127.0.0.1:13300")
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Fprintf(conn, "abc")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil || string(buf[0:n]) != "abc" {
			t.Error(err)
			return
		}
		forward.StopForward("abc")
	}
	{ //stop
		listenURL, _ := url.Parse("proxy://127.0.0.1:0")
		_, err := forward.StartForward("proxy", listenURL, "tcp://${HOST}")
		if err != nil {
			t.Error(err)
			return
		}
		forward.Stop()
	}
	{ //dial fail
		listenURL, _ := url.Parse("tcp://127.0.0.1:13300")
		_, err := forward.StartForward("abc", listenURL, "tcp://127.0.0.1:10")
		if err != nil {
			t.Error(err)
			return
		}
		conn, err := net.Dial("tcp", "127.0.0.1:13300")
		if err != nil {
			t.Error(err)
			return
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err == nil {
			t.Error(err)
			return
		}
		forward.StopForward("abc")
	}
	{ //error
		listenURL, _ := url.Parse("proxy://127.0.0.1:0")
		l, err := forward.StartForward("proxy", listenURL, "tcp://${HOST}")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = forward.StartForward("proxy", listenURL, "tcp://${HOST}")
		if err == nil {
			t.Error(err)
			return
		}
		forward.Stop()
		forward.procForward(l, "", nil, nil, "")
	}
	{ //dialer
		dialer := RouterPiperDialer{
			Next: xio.PiperDialerF(func(uri string, bufferSize int) (raw xio.Piper, err error) {
				return
			}),
		}
		dialer.DialPiper("", 0)
	}
}
