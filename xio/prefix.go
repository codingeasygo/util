package xio

import (
	"io"
	"net"
	"time"
)

// PrefixReader provide feature to prefix read data from Reader
type PrefixReader struct {
	io.Reader
	Prefix []byte
}

// NewPrefixReader will return new PrefixReader
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

// PreRead will read prefix size data before read loop start
func (p *PrefixReader) PreRead(size int) (data []byte, err error) {
	data = make([]byte, size)
	err = FullBuffer(p.Reader, data, uint32(size), nil)
	if err == nil {
		p.Prefix = data
	}
	return
}

func (p *PrefixReader) String() string {
	return RemoteAddr(p.Reader)
}

// PrefixReadWriteCloser is prefix read implement
type PrefixReadWriteCloser struct {
	io.ReadWriteCloser
	PrefixReader
}

// NewPrefixReadWriteCloser will return new PrefixReadWriteCloser
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

func (p *PrefixReadWriteCloser) String() string {
	return RemoteAddr(p.ReadWriteCloser)
}

// Network is net.Addr implement
func (p *PrefixReadWriteCloser) Network() string {
	return "prefix"
}

// LocalAddr returns the local network address.
func (p *PrefixReadWriteCloser) LocalAddr() net.Addr {
	if conn, ok := p.ReadWriteCloser.(net.Conn); ok {
		return conn.LocalAddr()
	}
	return p
}

// RemoteAddr returns the remote network address.
func (p *PrefixReadWriteCloser) RemoteAddr() net.Addr {
	if conn, ok := p.ReadWriteCloser.(net.Conn); ok {
		return conn.RemoteAddr()
	}
	return p
}

// SetDeadline sets the read and write deadlines associated
func (p *PrefixReadWriteCloser) SetDeadline(t time.Time) error {
	if conn, ok := p.ReadWriteCloser.(net.Conn); ok {
		return conn.SetDeadline(t)
	}
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
func (p *PrefixReadWriteCloser) SetReadDeadline(t time.Time) error {
	if conn, ok := p.ReadWriteCloser.(net.Conn); ok {
		return conn.SetReadDeadline(t)
	}
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
func (p *PrefixReadWriteCloser) SetWriteDeadline(t time.Time) error {
	if conn, ok := p.ReadWriteCloser.(net.Conn); ok {
		return conn.SetWriteDeadline(t)
	}
	return nil
}

// PrefixConn is net.Conn implement for prefix read data
type PrefixConn struct {
	net.Conn
	PrefixReader
}

// NewPrefixConn will return newPrefixConn
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

func (p *PrefixConn) String() string {
	return RemoteAddr(p.Conn)
}

// LocalAddr returns the local network address.
func (p *PrefixConn) LocalAddr() net.Addr {
	return p.Conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (p *PrefixConn) RemoteAddr() net.Addr {
	return p.Conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
func (p *PrefixConn) SetDeadline(t time.Time) error {
	return p.Conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
func (p *PrefixConn) SetReadDeadline(t time.Time) error {
	return p.Conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
func (p *PrefixConn) SetWriteDeadline(t time.Time) error {
	return p.Conn.SetWriteDeadline(t)
}
