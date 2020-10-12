package socks

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/codingeasygo/util/xio"
)

//Codable is interface for get current code
type Codable interface {
	Code() byte
}

//Server is an implementation of socks5 proxy
type Server struct {
	BufferSize int
	listners   map[net.Listener]string
	waiter     sync.WaitGroup
	Dialer     xio.PiperDialer
}

//NewServer will return new Server
func NewServer() (socks *Server) {
	socks = &Server{
		BufferSize: 32 * 1024,
		listners:   map[net.Listener]string{},
		waiter:     sync.WaitGroup{},
		Dialer:     xio.PiperDialerF(xio.DialNetPiper),
	}
	return
}

func (s *Server) loopAccept(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		go func(c net.Conn) {
			xerr := s.ProcConn(c)
			if xerr != xio.ErrAsyncRunning {
				c.Close()
			}
		}(conn)
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
func (s *Server) Start(addr string) (listener net.Listener, err error) {
	listener, err = net.Listen("tcp", addr)
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

//ProcConn will process connecton as socket protocol
func (s *Server) ProcConn(conn io.ReadWriteCloser) (err error) {
	DebugLog("Server proxy socks connection on %v from %v", xio.LocalAddr(conn), xio.RemoteAddr(conn))
	defer func() {
		if err != nil {
			DebugLog("Server proxy socks connection on %v from %v is done with %v", xio.LocalAddr(conn), xio.RemoteAddr(conn), err)
		}
		if err != xio.ErrAsyncRunning {
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
	switch buf[3] {
	case 0x01:
		err = xio.FullBuffer(conn, buf[5:], 5, nil)
		if err == nil {
			remote := fmt.Sprintf("%v.%v.%v.%v", buf[4], buf[5], buf[6], buf[7])
			port := uint16(buf[8])*256 + uint16(buf[9])
			uri = fmt.Sprintf("tcp://%v:%v", remote, port)
		}
	case 0x03:
		err = xio.FullBuffer(conn, buf[5:], uint32(buf[4]+2), nil)
		if err == nil {
			remote := string(buf[5 : buf[4]+5])
			port := uint16(buf[buf[4]+5])*256 + uint16(buf[buf[4]+6])
			uri = fmt.Sprintf("tcp://%v:%v", remote, port)
		}
	default:
		err = xio.FullBuffer(conn, buf[5:], uint32(buf[4]+2), nil)
		if err == nil {
			uri = string(buf[5 : buf[4]+5])
		}
	}
	// DebugLog("Server start dial to %v on %v", uri, conn.RemoteAddr())
	raw, err := s.Dialer.DialPiper(uri, s.BufferSize)
	if err != nil {
		buf[0], buf[1], buf[2], buf[3] = 0x05, 0x04, 0x00, 0x01
		buf[4], buf[5], buf[6], buf[7] = 0x00, 0x00, 0x00, 0x00
		buf[8], buf[9] = 0x00, 0x00
		if cerr, ok := err.(Codable); ok {
			buf[1] = cerr.Code()
		}
		conn.Write(buf[:10])
		// InfoLog("Server dial to %v on %v fail with %v", uri, conn.RemoteAddr(), err)
		return
	}
	buf[0], buf[1], buf[2], buf[3] = 0x05, 0x00, 0x00, 0x01
	buf[4], buf[5], buf[6], buf[7] = 0x00, 0x00, 0x00, 0x00
	buf[8], buf[9] = 0x00, 0x00
	_, err = conn.Write(buf[:10])
	if err != nil {
		raw.Close()
		return
	}
	err = raw.PipeConn(conn, uri)
	return
}

//Dial will dial connection by proxy server
func Dial(proxy, uri string) (conn net.Conn, err error) {
	conn, err = DialType(proxy, 0x03, uri)
	return
}

//DialType wil dial connection by proxy server and uri type
func DialType(proxy string, uriType byte, uri string) (conn net.Conn, err error) {
	proxy = strings.TrimPrefix(proxy, "socks5://")
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
	var host string
	var port int
	if uriType == 0x03 {
		var p string
		host, p, _ = net.SplitHostPort(uri)
		port, _ = strconv.Atoi(p)
	} else {
		host = uri
	}
	blen := len(host) + 7
	buf[0], buf[1], buf[2] = 0x05, 0x01, 0x00
	buf[3], buf[4] = uriType, byte(len(host))
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
