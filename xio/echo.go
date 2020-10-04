package xio

import (
	"io"
	"net"
	"time"
)

//EchoConn is net.Conn impl by os.Pipe
type EchoConn struct {
	piped *PipedChan
}

//NewEchoConn will return new echo connection
func NewEchoConn() (conn *EchoConn) {
	conn = &EchoConn{
		piped: NewPipedChan(),
	}
	return
}

func (e *EchoConn) Read(b []byte) (n int, err error) {
	n, err = e.piped.Read(b)
	return
}

func (e *EchoConn) Write(b []byte) (n int, err error) {
	n, err = e.piped.Write(b)
	return
}

// Close closes the connection.
func (e *EchoConn) Close() error {
	e.piped.Close()
	return nil
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

//EchoPiper is echo implement to Piper
type EchoPiper struct {
	BufferSize int
}

//NewEchoPiper will return new EchoPiper
func NewEchoPiper(bufferSize int) (piper *EchoPiper) {
	piper = &EchoPiper{BufferSize: bufferSize}
	return
}

//PipeConn will process connection by as echo
func (e *EchoPiper) PipeConn(conn net.Conn, target string) (err error) {
	_, err = io.CopyBuffer(conn, conn, make([]byte, e.BufferSize))
	return
}

//Close is empty implement for Piper
func (e *EchoPiper) Close() (err error) {
	return
}
