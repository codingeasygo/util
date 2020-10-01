package xio

import (
	"io"
	"net"
	"os"
	"time"
)

//PipeReadWriteCloser is pipe connection
type PipeReadWriteCloser struct {
	io.Reader
	io.Writer
	io.Closer
	Alias string
	side  *PipeReadWriteCloser
}

//Pipe will return new pipe connection.
func Pipe() (a, b *PipeReadWriteCloser, err error) {
	aReader, aWriter, err := os.Pipe()
	if err != nil {
		return
	}
	bReader, bWriter, err := os.Pipe()
	if err != nil {
		aWriter.Close()
		return
	}
	a = &PipeReadWriteCloser{
		Reader: aReader,
		Writer: bWriter,
		Closer: aWriter,
		Alias:  "piped",
	}
	b = &PipeReadWriteCloser{
		Reader: bReader,
		Writer: aWriter,
		Closer: bWriter,
		Alias:  "piped",
	}
	a.side = b
	b.side = a
	return
}

//Close will close reader/writer
func (p *PipeReadWriteCloser) Close() (err error) {
	err = p.Closer.Close()
	p.side.Closer.Close()
	return
}

func (p *PipeReadWriteCloser) String() string {
	return p.Alias
}

//PipedConn is an implementation of the net.Conn interface for piped two connection.
type PipedConn struct {
	*PipeReadWriteCloser
}

//CreatePipedConn will return two piped connection.
func CreatePipedConn() (a, b *PipedConn, err error) {
	basea, baseb, err := Pipe()
	if err == nil {
		a = &PipedConn{PipeReadWriteCloser: basea}
		b = &PipedConn{PipeReadWriteCloser: baseb}
	}
	return
}

//LocalAddr return self
func (p *PipedConn) LocalAddr() net.Addr {
	return p
}

//RemoteAddr return self
func (p *PipedConn) RemoteAddr() net.Addr {
	return p
}

//SetDeadline is empty
func (p *PipedConn) SetDeadline(t time.Time) error {
	return nil
}

//SetReadDeadline is empty
func (p *PipedConn) SetReadDeadline(t time.Time) error {
	return nil
}

//SetWriteDeadline is empty
func (p *PipedConn) SetWriteDeadline(t time.Time) error {
	return nil
}

//Network return "piped"
func (p *PipedConn) Network() string {
	return "piped"
}

func (p *PipedConn) String() string {
	return p.PipeReadWriteCloser.String()
}
