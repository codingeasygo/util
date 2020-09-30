package xio

import (
	"net"
	"os"
	"time"
)

//EchoConn is net.Conn impl by os.Pipe
type EchoConn struct {
	r, w *os.File
}

//NewEchoConn will return new echo connection
func NewEchoConn() (conn *EchoConn, err error) {
	conn = &EchoConn{}
	conn.r, conn.w, err = os.Pipe()
	return
}

func (e *EchoConn) Read(b []byte) (n int, err error) {
	n, err = e.r.Read(b)
	return
}

func (e *EchoConn) Write(b []byte) (n int, err error) {
	n, err = e.w.Write(b)
	return
}

// Close closes the connection.
func (e *EchoConn) Close() error {
	return e.r.Close()
}

// LocalAddr returns the local network address.
func (e *EchoConn) LocalAddr() net.Addr {
	return e
}

// RemoteAddr returns the remote network address.
func (e *EchoConn) RemoteAddr() net.Addr {
	return e
}

// SetDeadline for net.Conn
func (e *EchoConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline for net.Conn
func (e *EchoConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline for net.Conn
func (e *EchoConn) SetWriteDeadline(t time.Time) error {
	return nil
}

//Network is net.Addr impl
func (e *EchoConn) Network() string {
	return "echo"
}

func (e *EchoConn) String() string {
	return "echo"
}
