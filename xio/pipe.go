package xio

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PipedChan provoider Write buffer to Read implement by chan
type PipedChan struct {
	closed uint32
	having []byte
	locker *sync.Cond
}

// NewPipedChan will return new PipedChan
func NewPipedChan() (piped *PipedChan) {
	piped = &PipedChan{
		locker: sync.NewCond(&sync.Mutex{}),
	}
	return
}

func (p *PipedChan) Read(b []byte) (n int, err error) {
	p.locker.L.Lock()
	defer p.locker.L.Unlock()
	if len(p.having) < 1 {
		if p.closed > 0 {
			err = fmt.Errorf("closed")
			return
		}
		p.locker.Wait()
		if len(p.having) < 1 {
			err = fmt.Errorf("closed")
			return
		}
	}
	n = copy(b, p.having)
	p.having = p.having[n:]
	if len(p.having) < 1 {
		p.locker.Signal()
	}
	return
}

func (p *PipedChan) Write(b []byte) (n int, err error) {
	p.locker.L.Lock()
	defer p.locker.L.Unlock()
	if p.closed < 1 && len(p.having) > 0 {
		p.locker.Wait()
	}
	if p.closed > 0 {
		err = fmt.Errorf("closed")
		return
	}
	if len(p.having) > 0 {
		panic("having")
	}
	p.having = make([]byte, len(b))
	n = copy(p.having, b)
	p.locker.Signal()
	return
}

// Close will close piped channel
func (p *PipedChan) Close() (err error) {
	p.locker.L.Lock()
	defer p.locker.L.Unlock()
	if p.closed > 0 {
		err = fmt.Errorf("closed")
		return
	}
	p.closed = 1
	p.locker.Broadcast()
	return
}

// PipeReadWriteCloser is pipe connection
type PipeReadWriteCloser struct {
	Alias  string
	reader *PipedChan
	writer *PipedChan
	side   *PipeReadWriteCloser
}

// Pipe will return new pipe connection.
func Pipe() (a, b *PipeReadWriteCloser, err error) {
	piperA := NewPipedChan()
	piperB := NewPipedChan()

	a = &PipeReadWriteCloser{
		reader: piperA,
		writer: piperB,
		Alias:  "piped",
	}
	b = &PipeReadWriteCloser{
		reader: piperB,
		writer: piperA,
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

// Close will close reader/writer
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

// PipedConn is an implementation of the net.Conn interface for piped two connection.
type PipedConn struct {
	*PipeReadWriteCloser
}

// CreatePipedConn will return two piped connection.
func CreatePipedConn() (a, b *PipedConn, err error) {
	basea, baseb, err := Pipe()
	if err == nil {
		a = &PipedConn{PipeReadWriteCloser: basea}
		b = &PipedConn{PipeReadWriteCloser: baseb}
	}
	return
}

// LocalAddr return self
func (p *PipedConn) LocalAddr() net.Addr {
	return p
}

// RemoteAddr return self
func (p *PipedConn) RemoteAddr() net.Addr {
	return p
}

// SetDeadline is empty
func (p *PipedConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline is empty
func (p *PipedConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline is empty
func (p *PipedConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// Network return "piped"
func (p *PipedConn) Network() string {
	return "piped"
}

func (p *PipedConn) String() string {
	return p.PipeReadWriteCloser.String()
}

type PipedListener struct {
	queue chan *PipedConn
}

func NewPipedListener() (listener *PipedListener) {
	listener = &PipedListener{
		queue: make(chan *PipedConn, 8),
	}
	return
}

func (p *PipedListener) Dial() (conn net.Conn, err error) {
	conna, connb, err := CreatePipedConn()
	if err == nil {
		p.queue <- conna
		conn = connb
	}
	return
}

func (p *PipedListener) Accept() (conn net.Conn, err error) {
	c := <-p.queue
	if c == nil {
		err = fmt.Errorf("closed")
	} else {
		conn = c
	}
	return
}

func (p *PipedListener) Close() (err error) {
	p.queue <- nil
	return
}

func (p *PipedListener) Addr() net.Addr {
	return p
}

// Network return "piped"
func (p *PipedListener) Network() string {
	return "piped"
}

func (p *PipedListener) String() string {
	return fmt.Sprintf("%p", p)
}
