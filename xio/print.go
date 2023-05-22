package xio

import (
	"fmt"
	"io"
	"net"
	"time"
)

// PrintConn is net.Conn to print the transfter data
type PrintConn struct {
	Base io.ReadWriteCloser
	Name string
	Mode int
}

// NewPrintConn will create new PrintConn
func NewPrintConn(name string, base io.ReadWriteCloser) (conn *PrintConn) {
	conn = &PrintConn{Base: base, Name: name}
	return
}

func (p *PrintConn) Read(b []byte) (n int, err error) {
	n, err = p.Base.Read(b)
	if err == nil {
		switch p.Mode {
		case 0x10:
			fmt.Printf("%v Read %v bytes %v\n", p.Name, n, string(b[:n]))
		default:
			fmt.Printf("%v Read %v bytes % 02x\n", p.Name, n, b[:n])
		}
	} else {
		fmt.Printf("%v Read error %v\n", p.Name, err)
	}
	return
}

func (p *PrintConn) Write(b []byte) (n int, err error) {
	n, err = p.Base.Write(b)
	if err == nil {
		switch p.Mode {
		case 0x10:
			fmt.Printf("%v Write %v bytes %v\n", p.Name, n, string(b[:n]))
		default:
			fmt.Printf("%v Write %v bytes % 02x\n", p.Name, n, b[:n])
		}
	} else {
		fmt.Printf("%v Write error %v\n", p.Name, err)
	}
	return
}

// LocalAddr returns the local network address.
func (p *PrintConn) LocalAddr() net.Addr {
	if conn, ok := p.Base.(net.Conn); ok {
		return conn.LocalAddr()
	}
	return p
}

// RemoteAddr returns the remote network address.
func (p *PrintConn) RemoteAddr() net.Addr {
	if conn, ok := p.Base.(net.Conn); ok {
		return conn.RemoteAddr()
	}
	return p
}

// SetDeadline for net.Conn
func (p *PrintConn) SetDeadline(t time.Time) error {
	if conn, ok := p.Base.(net.Conn); ok {
		return conn.SetDeadline(t)
	}
	return nil
}

// SetReadDeadline for net.Conn
func (p *PrintConn) SetReadDeadline(t time.Time) error {
	if conn, ok := p.Base.(net.Conn); ok {
		return conn.SetReadDeadline(t)
	}
	return nil
}

// SetWriteDeadline for net.Conn
func (p *PrintConn) SetWriteDeadline(t time.Time) error {
	if conn, ok := p.Base.(net.Conn); ok {
		return conn.SetWriteDeadline(t)
	}
	return nil
}

// Network impl net.Addr
func (p *PrintConn) Network() string {
	return "print"
}

func (p *PrintConn) String() string {
	return fmt.Sprintf("%v", p.Name)
}

// Close will close base
func (p *PrintConn) Close() (err error) {
	err = p.Base.Close()
	fmt.Printf("%v Close %v\n", p.Name, err)
	return
}

type PrintPiper struct {
	Name string
	Raw  Piper
}

func NewPrintPiper(name string, raw Piper) (piper *PrintPiper) {
	piper = &PrintPiper{
		Name: name,
		Raw:  raw,
	}
	return
}

func (p *PrintPiper) PipeConn(conn io.ReadWriteCloser, target string) (err error) {
	err = p.Raw.PipeConn(NewPrintConn(p.Name, conn), target)
	return
}

func (p *PrintPiper) Close() (err error) {
	err = p.Raw.Close()
	return
}
