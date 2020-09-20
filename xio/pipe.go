package xio

import (
	"io"
	"os"
)

//PipeConn is pipe connection
type PipeConn struct {
	io.Reader
	io.Writer
	io.Closer
	side *PipeConn
}

//Pipe will return new pipe connection.
func Pipe() (a, b *PipeConn, err error) {
	aReader, aWriter, err := os.Pipe()
	if err != nil {
		return
	}
	bReader, bWriter, err := os.Pipe()
	if err != nil {
		aWriter.Close()
		return
	}
	a = &PipeConn{
		Reader: aReader,
		Writer: bWriter,
		Closer: aWriter,
	}
	b = &PipeConn{
		Reader: bReader,
		Writer: aWriter,
		Closer: bWriter,
	}
	a.side = b
	b.side = a
	return
}

//Close will close reader/writer
func (p *PipeConn) Close() (err error) {
	err = p.Closer.Close()
	p.side.Close()
	return
}
