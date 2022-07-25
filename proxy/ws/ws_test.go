package ws

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestProxy(t *testing.T) {
	testListener, _ := net.Listen("tcp", ":0")
	go func() {
		for {
			conn, err := testListener.Accept()
			if err != nil {
				break
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()
	testEcho := func(proxy, uri string) {
		conn, err := Dial(proxy, uri)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		fmt.Fprintf(conn, "abc")
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Error(err)
			return
		}
		if string(buffer[0:n]) != "abc" {
			t.Error("error")
			return
		}
	}
	server := NewServer()
	go func() {
		server.Run(":0")
	}()
	listener, _ := server.Start(":0")
	defer server.Stop()
	proxyServer := fmt.Sprintf("ws://%v", listener.Addr())
	testEcho(proxyServer, testListener.Addr().String())
	testEcho(proxyServer+"?abc=1", testListener.Addr().String())
	_, err := Dial(proxyServer, "")
	if err == nil {
		t.Error(err)
		return
	}
	_, err = Dial(proxyServer, "127.0.0.1:2")
	if err == nil {
		t.Error(err)
		return
	}
}
