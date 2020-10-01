package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/codingeasygo/util/xio"
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

//PreReadWriteCloser is prefix read implement
type PreReadWriteCloser struct {
	io.Reader
	io.Writer
	io.Closer
	readed []byte
}

func (p *PreReadWriteCloser) Read(b []byte) (n int, err error) {
	if p.readed != nil {
		n = copy(b, p.readed)
		p.readed = p.readed[n:]
		if len(p.readed) < 1 {
			p.readed = nil
		}
		if len(b) == n { //full
			return
		}
	}
	readed, err := p.Reader.Read(b[n:])
	n += readed
	return
}

//PreRead will read prefix size data before read loop start
func (p *PreReadWriteCloser) PreRead(size int) (data []byte, err error) {
	data = make([]byte, size)
	err = xio.FullBuffer(p.Reader, data, uint32(size), nil)
	if err == nil {
		p.readed = data
	}
	return
}

//PreConn is net.Conn implement for prefix read data
type PreConn struct {
	Pre PreReadWriteCloser
	net.Conn
}

func (p *PreConn) Read(b []byte) (n int, err error) {
	n, err = p.Pre.Read(b)
	return
}

//PreRead will read prefix size data before read loop start
func (p *PreConn) PreRead(size int) (data []byte, err error) {
	data, err = p.Pre.PreRead(size)
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
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			break
		}
		go b.ProcConn(conn)
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
	preConn := &PreConn{
		Conn: conn,
		Pre: PreReadWriteCloser{
			Reader: conn,
		},
	}
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
