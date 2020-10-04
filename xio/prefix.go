package xio

import (
	"io"
	"net"
)

//PrefixReader provide feature to prefix read data from Reader
type PrefixReader struct {
	io.Reader
	Prefix []byte
}

//NewPrefixReader will return new PrefixReader
func NewPrefixReader(base io.Reader) (prefix *PrefixReader) {
	prefix = &PrefixReader{Reader: base}
	return
}

func (p *PrefixReader) Read(b []byte) (n int, err error) {
	if p.Prefix != nil {
		n = copy(b, p.Prefix)
		p.Prefix = p.Prefix[n:]
		if len(p.Prefix) < 1 {
			p.Prefix = nil
		}
		return
	}
	readed, err := p.Reader.Read(b[n:])
	n += readed
	return
}

//PreRead will read prefix size data before read loop start
func (p *PrefixReader) PreRead(size int) (data []byte, err error) {
	data = make([]byte, size)
	err = FullBuffer(p.Reader, data, uint32(size), nil)
	if err == nil {
		p.Prefix = data
	}
	return
}

//PrefixReadWriteCloser is prefix read implement
type PrefixReadWriteCloser struct {
	io.ReadWriteCloser
	PrefixReader
}

//NewPrefixReadWriteCloser will return new PrefixReadWriteCloser
func NewPrefixReadWriteCloser(base io.ReadWriteCloser) (prefix *PrefixReadWriteCloser) {
	prefix = &PrefixReadWriteCloser{}
	prefix.ReadWriteCloser = base
	prefix.PrefixReader.Reader = base
	return
}

func (p *PrefixReadWriteCloser) Read(b []byte) (n int, err error) {
	n, err = p.PrefixReader.Read(b)
	return
}

//PrefixConn is net.Conn implement for prefix read data
type PrefixConn struct {
	net.Conn
	PrefixReader
}

//NewPrefixConn will return newPrefixConn
func NewPrefixConn(conn net.Conn) (prefix *PrefixConn) {
	prefix = &PrefixConn{}
	prefix.Conn = conn
	prefix.PrefixReader.Reader = conn
	return
}

func (p *PrefixConn) Read(b []byte) (n int, err error) {
	n, err = p.PrefixReader.Read(b)
	return
}
