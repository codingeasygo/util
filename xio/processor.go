package xio

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

//Processor is interface for process connection
type Processor interface {
	ProcConn(conn io.ReadWriteCloser) (err error)
}

//ProcessorF is func to implement Processor
type ProcessorF func(conn io.ReadWriteCloser) (err error)

//ProcConn is process connection by func
func (p ProcessorF) ProcConn(conn io.ReadWriteCloser) (err error) {
	err = p(conn)
	return
}

//ErrAsyncRunning is error for async running on PipeConn
var ErrAsyncRunning = fmt.Errorf("asynced")

//Piper is interface for process pipe connection
type Piper interface {
	PipeConn(conn io.ReadWriteCloser, target string) (err error)
	Close() (err error)
}

//PiperF is func to implement Piper
type PiperF func(conn io.ReadWriteCloser, target string) (err error)

func (p PiperF) PipeConn(conn io.ReadWriteCloser, target string) (err error) {
	err = p(conn, target)
	return
}

func (p PiperF) Close() (err error) {
	return
}

//PiperDialer is interface for implement piper dialer
type PiperDialer interface {
	DialPiper(uri string, bufferSize int) (raw Piper, err error)
}

//PiperDialerF is func to implement PiperDialer
type PiperDialerF func(uri string, bufferSize int) (raw Piper, err error)

//DialPiper will dial one piper by uri
func (p PiperDialerF) DialPiper(uri string, bufferSize int) (raw Piper, err error) {
	raw, err = p(uri, bufferSize)
	return
}

//NetPiper is Piper implement by net.Dial
type NetPiper struct {
	net.Conn
	CopyPiper
}

//DialNetPiper will return new NetPiper by net.Dial
func DialNetPiper(uri string, bufferSize int) (piper Piper, err error) {
	var network, address string
	parts := strings.SplitN(uri, "://", 2)
	if len(parts) < 2 {
		network = "tcp"
		address = parts[0]
	} else {
		network = parts[0]
		address = parts[1]
	}
	conn, err := net.Dial(network, address)
	if err == nil {
		piper = &NetPiper{
			Conn: conn,
			CopyPiper: CopyPiper{
				ReadWriteCloser: conn,
				BufferSize:      bufferSize,
			},
		}
	}
	return
}

//CopyPiper is Piper implement by copy
type CopyPiper struct {
	io.ReadWriteCloser
	BufferSize int
	XX         string
}

//NewCopyPiper will return new CopyPiper
func NewCopyPiper(raw io.ReadWriteCloser, bufferSize int) (piper *CopyPiper) {
	piper = &CopyPiper{ReadWriteCloser: raw, BufferSize: bufferSize}
	return
}

//PipeConn will pipe connection to raw
func (c *CopyPiper) PipeConn(conn io.ReadWriteCloser, target string) (err error) {
	wc := make(chan int, 1)
	go func() {
		var readErr error
		if to, ok := conn.(io.WriterTo); ok {
			_, readErr = to.WriteTo(c.ReadWriteCloser)
		} else if from, ok := c.ReadWriteCloser.(io.ReaderFrom); ok {
			_, readErr = from.ReadFrom(conn)
		} else {
			_, readErr = io.CopyBuffer(c.ReadWriteCloser, conn, make([]byte, c.BufferSize))
		}
		c.ReadWriteCloser.Close()
		wc <- 1
		if err == nil {
			err = readErr
		}
	}()
	{
		var writeErr error
		if from, ok := conn.(io.ReaderFrom); ok {
			_, writeErr = from.ReadFrom(c.ReadWriteCloser)
		} else if to, ok := c.ReadWriteCloser.(io.WriterTo); ok {
			_, writeErr = to.WriteTo(conn)
		} else {
			_, writeErr = io.CopyBuffer(conn, c.ReadWriteCloser, make([]byte, c.BufferSize))
		}
		conn.Close()
		if err == nil {
			err = writeErr
		}
	}
	<-wc
	close(wc)
	return
}

//ByteDistributeProcessor is distribute processor by prefix read first byte
type ByteDistributeProcessor struct {
	Next      map[byte]Processor
	conns     map[string]io.ReadWriteCloser
	listeners map[string]net.Listener
	locker    sync.RWMutex
}

//NewByteDistributeProcessor will return new processor
func NewByteDistributeProcessor() (processor *ByteDistributeProcessor) {
	processor = &ByteDistributeProcessor{
		Next:      map[byte]Processor{},
		conns:     map[string]io.ReadWriteCloser{},
		listeners: map[string]net.Listener{},
		locker:    sync.RWMutex{},
	}
	return
}

//AddProcessor will add processor by mode
func (b *ByteDistributeProcessor) AddProcessor(m byte, procesor Processor) {
	b.Next[m] = procesor
}

//RemoveProcessor will remove processor by mode
func (b *ByteDistributeProcessor) RemoveProcessor(m byte) {
	delete(b.Next, m)
}

//ProcAccept will loop accept net.Conn and async call ProcConn
func (b *ByteDistributeProcessor) ProcAccept(listener net.Listener) (err error) {
	b.locker.Lock()
	b.listeners[fmt.Sprintf("%p", listener)] = listener
	b.locker.Unlock()
	defer func() {
		b.locker.Lock()
		delete(b.listeners, fmt.Sprintf("%p", listener))
		b.locker.Unlock()
	}()
	procConn := func(c net.Conn) {
		xerr := b.ProcConn(c)
		if xerr != ErrAsyncRunning {
			c.Close()
		}
	}
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			break
		}
		go procConn(conn)
	}
	return
}

//ProcConn will process connection by prefix reader and distribute to next processor
func (b *ByteDistributeProcessor) ProcConn(conn io.ReadWriteCloser) (err error) {
	b.locker.Lock()
	b.conns[fmt.Sprintf("%p", conn)] = conn
	b.locker.Unlock()
	defer func() {
		b.locker.Lock()
		delete(b.conns, fmt.Sprintf("%p", conn))
		b.locker.Unlock()
	}()
	preConn := NewPrefixReadWriteCloser(conn)
	pre, err := preConn.PreRead(1)
	if err != nil {
		return
	}
	processor := b.Next[pre[0]]
	if processor == nil {
		processor = b.Next['*']
	}
	if processor == nil {
		err = fmt.Errorf("processor is not exist by %v", pre[0])
		return
	}
	err = processor.ProcConn(preConn)
	return
}

//Close will close all listener and connection
func (b *ByteDistributeProcessor) Close() (err error) {
	b.locker.Lock()
	for _, conn := range b.conns {
		conn.Close()
	}
	for _, listener := range b.listeners {
		listener.Close()
	}
	b.locker.Unlock()
	return
}
