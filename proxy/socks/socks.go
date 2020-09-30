package socks

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/codingeasygo/util/xio"
)

//Codable is interface for get current code
type Codable interface {
	Code() byte
}

//PendingConn is an implementation of io.ReadWriteCloser
type PendingConn struct {
	Raw     io.ReadWriteCloser
	having  []byte
	pending uint32
	wc      chan int
}

//NewPendingConn will return new endingConn
func NewPendingConn(raw io.ReadWriteCloser, having []byte) (conn *PendingConn) {
	conn = &PendingConn{
		Raw:     raw,
		having:  having,
		pending: 1,
		wc:      make(chan int),
	}
	return
}

//Start pending connection
func (p *PendingConn) Start() {
	if atomic.CompareAndSwapUint32(&p.pending, 1, 0) {
		close(p.wc)
	}
}

func (p *PendingConn) Write(b []byte) (n int, err error) {
	if p.pending == 1 {
		<-p.wc
	}
	n, err = p.Raw.Write(b)
	return
}

func (p *PendingConn) Read(b []byte) (n int, err error) {
	if p.pending == 1 {
		<-p.wc
	}
	if len(p.having) > 0 {
		n = copy(b, p.having)
		p.having = p.having[0:n]
	} else {
		n, err = p.Raw.Read(b)
	}
	return
}

//Close pending connection.
func (p *PendingConn) Close() (err error) {
	if atomic.CompareAndSwapUint32(&p.pending, 1, 0) {
		close(p.wc)
	}
	err = p.Raw.Close()
	return
}

const (
	//URITypeNormal is normal URI type
	URITypeNormal = 0
	//URITypeBS is b
	URITypeBS = 1
)

//Server is an implementation of socks5 proxy
type Server struct {
	listners map[net.Listener]string
	waiter   sync.WaitGroup
	Dialer   func(utype int, uri string, conn io.ReadWriteCloser) (raw io.ReadWriteCloser, err error)
}

//NewServer will return new Server
func NewServer() (socks *Server) {
	socks = &Server{
		listners: map[net.Listener]string{},
		waiter:   sync.WaitGroup{},
		Dialer: func(utype int, uri string, conn io.ReadWriteCloser) (raw io.ReadWriteCloser, err error) {
			raw, err = net.Dial("tcp", uri)
			return
		},
	}
	return
}

func (s *Server) loopAccept(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		go s.procConn(conn)
	}
	s.waiter.Done()
}

//Run will listen tcp on address and sync accept to ProcConn
func (s *Server) Run(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err == nil {
		s.listners[listener] = addr
		InfoLog("Server listen http proxy on %v", addr)
		s.waiter.Add(1)
		s.loopAccept(listener)
	}
	return
}

//Start proxy listener
func (s *Server) Start(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err == nil {
		s.listners[listener] = addr
		InfoLog("Server listen socks5 proxy on %v", addr)
		s.waiter.Add(1)
		go s.loopAccept(listener)
	}
	return
}

//Stop will stop listener and wait loop stop
func (s *Server) Stop() (err error) {
	for listener, addr := range s.listners {
		err = listener.Close()
		delete(s.listners, listener)
		InfoLog("Server socks5 proxy listener on %v is stopped by %v", addr, err)
	}
	s.waiter.Wait()
	return
}

func (s *Server) procConn(conn net.Conn) (err error) {
	var raw io.ReadWriteCloser
	DebugLog("Server proxy socks connection from %v", conn.RemoteAddr())
	defer func() {
		if err != nil {
			DebugLog("Server proxy socks connection from %v is done with %v", conn.RemoteAddr(), err)
		}
		if raw != nil {
			conn.Close()
		}
	}()
	buf := make([]byte, 1024*64)
	//
	//Procedure method
	err = xio.FullBuffer(conn, buf, 2, nil)
	if err != nil {
		return
	}
	if buf[0] != 0x05 {
		err = fmt.Errorf("only ver 0x05 is supported, but %x", buf[0])
		return
	}
	err = xio.FullBuffer(conn, buf[2:], uint32(buf[1]), nil)
	if err != nil {
		return
	}
	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		return
	}
	//
	//Procedure request
	err = xio.FullBuffer(conn, buf, 5, nil)
	if err != nil {
		return
	}
	if buf[0] != 0x05 {
		err = fmt.Errorf("only ver 0x05 is supported, but %x", buf[0])
		return
	}
	var uri string
	var utype int
	switch buf[3] {
	case 0x01:
		err = xio.FullBuffer(conn, buf[5:], 5, nil)
		if err == nil {
			remote := fmt.Sprintf("%v.%v.%v.%v", buf[4], buf[5], buf[6], buf[7])
			port := uint16(buf[8])*256 + uint16(buf[9])
			uri = fmt.Sprintf("%v:%v", remote, port)
			utype = URITypeNormal
		}
	case 0x03:
		err = xio.FullBuffer(conn, buf[5:], uint32(buf[4]+2), nil)
		if err == nil {
			remote := string(buf[5 : buf[4]+5])
			port := uint16(buf[buf[4]+5])*256 + uint16(buf[buf[4]+6])
			uri = fmt.Sprintf("%v:%v", remote, port)
			utype = URITypeNormal
		}
	case 0x13:
		err = xio.FullBuffer(conn, buf[5:], uint32(buf[4]+2), nil)
		if err == nil {
			uri = string(buf[5 : buf[4]+5])
			utype = URITypeBS
		}
	default:
		err = fmt.Errorf("ATYP %v is not supported", buf[3])
		return
	}
	// DebugLog("Server start dial to %v on %v", uri, conn.RemoteAddr())
	pending := NewPendingConn(conn, nil)
	raw, err = s.Dialer(utype, uri, pending)
	if err != nil {
		buf[0], buf[1], buf[2], buf[3] = 0x05, 0x04, 0x00, 0x01
		buf[4], buf[5], buf[6], buf[7] = 0x00, 0x00, 0x00, 0x00
		buf[8], buf[9] = 0x00, 0x00
		if cerr, ok := err.(Codable); ok {
			buf[1] = cerr.Code()
		}
		conn.Write(buf[:10])
		// InfoLog("Server dial to %v on %v fail with %v", uri, conn.RemoteAddr(), err)
		pending.Close()
		return
	}
	buf[0], buf[1], buf[2], buf[3] = 0x05, 0x00, 0x00, 0x01
	buf[4], buf[5], buf[6], buf[7] = 0x00, 0x00, 0x00, 0x00
	buf[8], buf[9] = 0x00, 0x00
	_, err = conn.Write(buf[:10])
	if err != nil {
		pending.Close()
		return
	}
	pending.Start()
	if raw != nil {
		go func() {
			io.Copy(raw, pending)
			raw.Close()
		}()
		io.Copy(pending, raw)
	}
	return
}

//Dial will dial connection by proxy server
func Dial(proxy, uri string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", proxy)
	if err != nil {
		return
	}
	conn.Write([]byte{0x05, 0x01, 0x00})
	buf := make([]byte, 1024*64)
	err = xio.FullBuffer(conn, buf, 2, nil)
	if err != nil {
		conn.Close()
		return
	}
	if buf[0] != 0x05 || buf[1] != 0x00 {
		err = fmt.Errorf("unsupported %x", buf)
		conn.Close()
		return
	}
	host, p, _ := net.SplitHostPort(uri)
	port, _ := strconv.Atoi(p)
	blen := len(host) + 7
	buf[0], buf[1], buf[2] = 0x05, 0x01, 0x00
	buf[3], buf[4] = 0x03, byte(len(host))
	copy(buf[5:], []byte(host))
	buf[blen-2] = byte(port / 256)
	buf[blen-1] = byte(port % 256)
	conn.Write(buf[:blen])
	err = xio.FullBuffer(conn, buf, 5, nil)
	if err != nil {
		conn.Close()
		return
	}
	switch buf[3] {
	case 0x01:
		err = xio.FullBuffer(conn, buf[5:], 5, nil)
	case 0x03:
		err = xio.FullBuffer(conn, buf[5:], uint32(buf[4])+2, nil)
	case 0x04:
		err = xio.FullBuffer(conn, buf[5:], 17, nil)
	default:
		err = fmt.Errorf("reply address type is not supported:%v", buf[3])
	}
	if err != nil {
		conn.Close()
		return
	}
	if buf[1] != 0x00 {
		conn.Close()
		err = fmt.Errorf("response code(%x)", buf[1])
		return
	}
	return
}
