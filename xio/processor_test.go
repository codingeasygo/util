package xio

import (
	"fmt"
	"io"
	"net"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestNetPiper(t *testing.T) {
	listener, _ := net.Listen("tcp", ":0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			go io.Copy(conn, conn)
		}
	}()
	dialer := PiperDialerF(DialNetPiper)
	{
		uri := listener.Addr().String()
		piper, err := dialer.DialPiper(uri, 1024)
		if err != nil {
			t.Error(err)
			return
		}
		conna, connb, _ := CreatePipedConn()
		go piper.PipeConn(connb, uri)
		fmt.Fprintf(conna, "abc")
		buffer := make([]byte, 1024)
		err = FullBuffer(conna, buffer, 3, nil)
		if err != nil || string(buffer[0:3]) != "abc" {
			t.Error(err)
			return
		}
		conna.Close()
		time.Sleep(10 * time.Millisecond)
	}
	{
		uri := "tcp://" + listener.Addr().String()
		piper, err := dialer.DialPiper(uri, 1024)
		if err != nil {
			t.Error(err)
			return
		}
		conna, connb, _ := CreatePipedConn()
		go piper.PipeConn(connb, uri)
		fmt.Fprintf(conna, "abc")
		buffer := make([]byte, 1024)
		err = FullBuffer(conna, buffer, 3, nil)
		if err != nil || string(buffer[0:3]) != "abc" {
			t.Error(err)
			return
		}
		conna.Close()
		time.Sleep(10 * time.Millisecond)
	}
}

func TestByteDistribute(t *testing.T) {
	accept := make(chan net.Conn, 1)
	processor := NewByteDistributeProcessor()
	go processor.ProcAccept(ListenerF(func() (conn net.Conn, err error) {
		conn = <-accept
		if conn == nil {
			err = fmt.Errorf("closed")
		}
		return
	}))
	processor.AddProcessor('A', ProcessorF(func(conn net.Conn) (err error) {
		buf := make([]byte, 1024)
		for {
			err = FullBuffer(conn, buf, 4, nil)
			if err != nil {
				break
			}
			fmt.Fprintf(conn, "A:%v", string(buf[1:4]))
		}
		return
	}))
	processor.AddProcessor('*', ProcessorF(func(conn net.Conn) (err error) {
		buf := make([]byte, 1024)
		for {
			err = FullBuffer(conn, buf, 4, nil)
			if err != nil {
				break
			}
			fmt.Fprintf(conn, "*:%v", string(buf[1:4]))
		}
		return
	}))
	buf := make([]byte, 1024)
	//
	//A
	conna, connb, _ := CreatePipedConn()
	accept <- connb
	fmt.Fprintf(conna, "A123")
	readed, err := conna.Read(buf)
	if err != nil || string(buf[0:readed]) != "A:123" {
		t.Errorf("err:%v,%v", err, string(buf[0:readed]))
		return
	}
	conna.Close()
	//
	//B
	conna, connb, _ = CreatePipedConn()
	accept <- connb
	fmt.Fprintf(conna, "B123")
	readed, err = conna.Read(buf)
	if err != nil || string(buf[0:readed]) != "*:123" {
		t.Error(err)
		return
	}
	//
	//B
	processor.RemoveProcessor('*')
	conna, connb, _ = CreatePipedConn()
	accept <- connb
	fmt.Fprintf(conna, "B123")
	readed, err = conna.Read(buf)
	if err == nil {
		t.Error(err)
		return
	}
	// conna.Close()
	//
	//error
	conna, connb, _ = CreatePipedConn()
	accept <- connb
	conna.Close()
	//
	//
	processor.Close()
	accept <- nil
	time.Sleep(10 * time.Millisecond)
}
