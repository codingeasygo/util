package xio

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

type printConnTest struct {
	data chan []byte
	err  error
}

func (p *printConnTest) Read(b []byte) (n int, err error) {
	if p.err == nil {
		data := <-p.data
		n = copy(b, data)
	}
	err = p.err
	return
}

func (p *printConnTest) Write(b []byte) (n int, err error) {
	if p.err == nil {
		p.data <- b
		n = len(b)
	}
	err = p.err
	return
}

// LocalAddr returns the local network address.
func (p *printConnTest) LocalAddr() net.Addr {
	return p
}

// RemoteAddr returns the remote network address.
func (p *printConnTest) RemoteAddr() net.Addr {
	return p
}

// SetDeadline for net.Conn
func (p *printConnTest) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline for net.Conn
func (p *printConnTest) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline for net.Conn
func (p *printConnTest) SetWriteDeadline(t time.Time) error {
	return nil
}

//Network impl net.Addr
func (p *printConnTest) Network() string {
	return "print"
}

func (p *printConnTest) String() string {
	return "test"
}

func (p *printConnTest) Close() (err error) {
	return
}

func TestPrint(t *testing.T) {
	test := &printConnTest{
		data: make(chan []byte, 1024),
	}
	print := NewPrintConn("testing", test)
	//
	print.Mode = 0
	print.Write([]byte("abc"))
	print.Read(make([]byte, 1024))
	//
	print.Mode = 0x10
	print.Write([]byte("abc"))
	print.Read(make([]byte, 1024))
	//
	test.err = fmt.Errorf("closed")
	print.Write([]byte("abc"))
	print.Read(make([]byte, 1024))
	//
	print.Close()
	//
	print.SetDeadline(time.Now())
	print.SetReadDeadline(time.Now())
	print.SetWriteDeadline(time.Now())
	print.LocalAddr()
	print.RemoteAddr()
	print.Network()
	fmt.Println(print.String())
	//
	test1 := bytes.NewBuffer(nil)
	print1 := NewPrintConn("testing", NewCombinedReadWriteCloser(test1, test1, nil))
	print1.SetDeadline(time.Now())
	print1.SetReadDeadline(time.Now())
	print1.SetWriteDeadline(time.Now())
	print1.LocalAddr()
	print1.RemoteAddr()
	//
	print2 := NewPrintPiper("testing", PiperF(func(conn io.ReadWriteCloser, target string) (err error) {
		err = fmt.Errorf("error")
		return
	}))
	print2.PipeConn(nil, "")
	print2.Close()
}
