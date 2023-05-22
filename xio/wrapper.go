package xio

import (
	"fmt"
	"io"
	"net"
	"time"
)

// ReaderF is wrapper for io.Reader
type ReaderF func(p []byte) (n int, err error)

func (r ReaderF) Read(p []byte) (n int, err error) {
	n, err = r(p)
	return
}

// WriterF is wrapper for io.Writer
type WriterF func(p []byte) (n int, err error)

func (w WriterF) Write(p []byte) (n int, err error) {
	n, err = w(p)
	return
}

// CloserF is wrapper for io.Closer
type CloserF func() (err error)

// Close will close by func
func (c CloserF) Close() (err error) {
	err = c()
	return
}

// ConnWrapper is wrapper for net.Conn by ReadWriteCloser
type ConnWrapper struct {
	io.ReadWriteCloser
}

// NewConnWrapper will create new ConnWrapper
func NewConnWrapper(base io.ReadWriteCloser) (wrapper *ConnWrapper) {
	return &ConnWrapper{ReadWriteCloser: base}
}

// Network impl net.Addr
func (c *ConnWrapper) Network() string {
	return "wrapper"
}

func (c *ConnWrapper) String() string {
	return fmt.Sprintf("%v", c.ReadWriteCloser)
}

// LocalAddr return then local network address
func (c *ConnWrapper) LocalAddr() net.Addr {
	return c
}

// RemoteAddr returns the remote network address.
func (c *ConnWrapper) RemoteAddr() net.Addr {
	return c
}

// SetDeadline impl net.Conn do nothing
func (c *ConnWrapper) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline impl net.Conn do nothing
func (c *ConnWrapper) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline impl net.Conn do nothing
func (c *ConnWrapper) SetWriteDeadline(t time.Time) error {
	return nil
}
