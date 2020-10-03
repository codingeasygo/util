package xio

import (
	"net"
	"os"
	"time"
)

//PipeReadWriteCloser is pipe connection
type PipeReadWriteCloser struct {
	Alias  string
	reader *os.File
	writer *os.File
	side   *PipeReadWriteCloser
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
		reader: aReader,
		writer: bWriter,
		Alias:  "piped",
	}
	b = &PipeReadWriteCloser{
		reader: bReader,
		writer: aWriter,
		Alias:  "piped",
	}
	a.side = b
	b.side = a
	return
}

func (p *PipeReadWriteCloser) Read(b []byte) (n int, err error) {
	n, err = p.reader.Read(b)
	return
}

func (p *PipeReadWriteCloser) Write(b []byte) (n int, err error) {
	n, err = p.writer.Write(b)
	return
}

//Close will close reader/writer
func (p *PipeReadWriteCloser) Close() (err error) {
	p.reader.Close()
	p.writer.Close()
	p.side.reader.Close()
	p.side.writer.Close()
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
