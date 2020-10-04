package http

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/codingeasygo/util/xio"
)

//Server is http proxy server
type Server struct {
	BufferSize int
	listners   map[net.Listener]string
	waiter     sync.WaitGroup
	Dialer     xio.PiperDialer
}

//NewServer will return new server
func NewServer() (proxy *Server) {
	proxy = &Server{
		BufferSize: 32 * 1024,
		listners:   map[net.Listener]string{},
		waiter:     sync.WaitGroup{},
		Dialer:     xio.PiperDialerF(xio.DialNetPiper),
	}
	return
}

//Run will listen tcp on address and accept to ProcConn
func (s *Server) loopAccept(l net.Listener) (err error) {
	var conn net.Conn
	for {
		conn, err = l.Accept()
		if err != nil {
			break
		}
		go func() {
			xerr := s.ProcConn(conn)
			if xerr != xio.ErrAsyncRunning {
				conn.Close()
			}
		}()
	}
	s.waiter.Done()
	return
}

//Run will listen tcp on address and sync accept to ProcConn
func (s *Server) Run(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err == nil {
		s.listners[listener] = addr
		InfoLog("Server listen http proxy on %v", addr)
		s.waiter.Add(1)
		err = s.loopAccept(listener)
	}
	return
}

//Start will listen tcp on address and async accept to ProcConn
func (s *Server) Start(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err == nil {
		s.listners[listener] = addr
		InfoLog("Server listen http proxy on %v", addr)
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
		InfoLog("Server http proxy listener on %v is stopped by %v", addr, err)
	}
	s.waiter.Wait()
	return
}

//ProcConn will processs net connect as http proxy
func (s *Server) ProcConn(conn net.Conn) (err error) {
	DebugLog("Server proxy http connection from %v", conn.RemoteAddr())
	defer func() {
		if err != nil {
			DebugLog("Server proxy http connection from %v is done with %v", conn.RemoteAddr(), err)
		}
	}()
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		return
	}
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("Proxy-Connection")
	resp := &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
	}
	resp.Header.Add("Proxy-Agent", "test/v1.0.0")
	var raw xio.Piper
	var uri string
	if req.Method == "CONNECT" {
		uri = "tcp://" + req.RequestURI
		raw, err = s.Dialer.DialPiper(uri, s.BufferSize)
		if err != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.Write(conn)
			return
		}
		resp.StatusCode = http.StatusOK
		resp.Status = "Connection established"
		resp.Write(conn)
	} else {
		host := req.Host
		if _, port, _ := net.SplitHostPort(host); port == "" {
			host += ":80"
		}
		uri = "tcp://" + host
		raw, err = s.Dialer.DialPiper(uri, s.BufferSize)
		if err != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.Write(conn)
			return
		}
		buffer := bytes.NewBuffer(nil)
		req.Write(buffer)
		prefix := xio.NewPrefixConn(conn)
		prefix.Prefix = buffer.Bytes()
		conn = prefix
	}
	err = raw.PipeConn(conn, uri)
	return
}

//Dial will dial uri by proxy server
func Dial(proxy, uri string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", proxy)
	if err != nil {
		return
	}
	address := uri
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	req, err := http.NewRequest("CONNECT", address, nil)
	if err != nil {
		return
	}
	req.Write(conn)
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, req)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("proxy response %v", resp.StatusCode)
	}
	return
}
