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
	ProcConn(conn net.Conn) (err error)
}

//ProcessorF is func to implement Processor
type ProcessorF func(conn net.Conn) (err error)

//ProcConn is process connection by func
func (p ProcessorF) ProcConn(conn net.Conn) (err error) {
	err = p(conn)
	return
}

//ErrAsyncRunning is error for async running on PipeConn
var ErrAsyncRunning = fmt.Errorf("asynced")

//Piper is interface for process pipe connection
type Piper interface {
	PipeConn(conn net.Conn, target string) (err error)
	Close() (err error)
}

//PiperF is func to implement Piper
type PiperF func(conn net.Conn, target string) (err error)

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
	BufferSize int
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
		piper = &NetPiper{Conn: conn, BufferSize: bufferSize}
	}
	return
}

//PipeConn will pipe conn to n.Conn by copy
func (n *NetPiper) PipeConn(conn net.Conn, target string) (err error) {
	wc := make(chan int, 1)
	go func() {
		io.CopyBuffer(n.Conn, conn, make([]byte, n.BufferSize))
		n.Conn.Close()
		wc <- 1
	}()
	_, err = io.CopyBuffer(conn, n.Conn, make([]byte, n.BufferSize))
	n.Close()
	<-wc
	return
}

//ByteDistributeProcessor is distribute processor by prefix read first byte
type ByteDistributeProcessor struct {
	Next      map[byte]Processor
	conns     map[string]net.Conn
	listeners map[string]net.Listener
	locker    sync.RWMutex
}

//NewByteDistributeProcessor will return new processor
func NewByteDistributeProcessor() (processor *ByteDistributeProcessor) {
	processor = &ByteDistributeProcessor{
		Next:      map[byte]Processor{},
		conns:     map[string]net.Conn{},
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
func (b *ByteDistributeProcessor) ProcConn(conn net.Conn) (err error) {
	b.locker.Lock()
	b.conns[fmt.Sprintf("%p", conn)] = conn
	b.locker.Unlock()
	defer func() {
		b.locker.Lock()
		delete(b.conns, fmt.Sprintf("%p", conn))
		b.locker.Unlock()
	}()
	preConn := NewPrefixConn(conn)
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