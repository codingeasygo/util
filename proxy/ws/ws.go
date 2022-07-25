package ws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/codingeasygo/util/xio"
	"github.com/codingeasygo/util/xnet"
	"golang.org/x/net/websocket"
)

type ContextKey string

type Server struct {
	*websocket.Server
	BufferSize int
	Dialer     xio.PiperDialer
	waiter     sync.WaitGroup
	listners   map[net.Listener]string
}

func NewServer() (server *Server) {
	server = &Server{
		BufferSize: 32 * 1024,
		Dialer:     xio.PiperDialerF(xio.DialNetPiper),
		waiter:     sync.WaitGroup{},
		listners:   map[net.Listener]string{},
	}
	server.Server = &websocket.Server{Handler: server.handler}
	return
}

//Run will listen tcp on address and accept to ProcConn
func (s *Server) loopAccept(l net.Listener) (err error) {
	defer s.waiter.Done()
	http.Serve(l, s)
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
func (s *Server) Start(addr string) (listener net.Listener, err error) {
	listener, err = net.Listen("tcp", addr)
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

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	uri := req.Form.Get("_uri")
	if len(uri) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "_uri is required")
		return
	}
	raw, err := s.Dialer.DialPiper(uri, s.BufferSize)
	if err != nil {
		InfoLog("Server dial to %v fail with %v", uri, err)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "%v", err)
		return
	}
	newReq := context.WithValue(req.Context(), ContextKey("upstream"), []interface{}{raw, uri})
	s.Server.ServeHTTP(w, req.WithContext(newReq))
}

func (s *Server) handler(conn *websocket.Conn) {
	defer conn.Close()
	req := conn.Request()
	upstream := req.Context().Value(ContextKey("upstream")).([]interface{})
	raw, uri := upstream[0].(xio.Piper), upstream[1].(string)
	DebugLog("Server start forward %v to %v", req.RemoteAddr, uri)
	err := raw.PipeConn(conn, uri)
	DebugLog("Server forward %v to %v is done with %v", req.RemoteAddr, uri, err)
}

//Dial will dial connection by proxy server
func Dial(proxy, uri string) (conn net.Conn, err error) {
	dialer := xnet.NewWebsocketDialer()
	targetURI := proxy
	if strings.Contains(proxy, "?") {
		targetURI += fmt.Sprintf("&_uri=%v", url.QueryEscape(uri))
	} else {
		targetURI += fmt.Sprintf("?_uri=%v", url.QueryEscape(uri))
	}
	raw, err := dialer.Dial(targetURI)
	if err == nil {
		conn = raw.(net.Conn)
	}
	return
}
