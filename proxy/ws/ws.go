package ws

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/codingeasygo/util/xio"
	"github.com/codingeasygo/util/xnet"
	"golang.org/x/net/websocket"
)

type Server struct {
	*websocket.Server
	BufferSize int
	Dialer     xio.PiperDialer
}

func NewServer() (server *Server) {
	server = &Server{
		BufferSize: 32 * 1024,
		Dialer:     xio.PiperDialerF(xio.DialNetPiper),
	}
	server.Server = &websocket.Server{Handler: server.handleWS}
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
	s.Server.ServeHTTP(w, req)
}

func (s *Server) handleWS(conn *websocket.Conn) {
	defer conn.Close()
	req := conn.Request()
	uri := conn.Request().Form.Get("_uri")
	raw, err := s.Dialer.DialPiper(uri, s.BufferSize)
	if err != nil {
		InfoLog("Server dial to %v fail with %v", uri, err)
		return
	}
	DebugLog("Server start forward %v to %v", req.RemoteAddr, uri)
	err = raw.PipeConn(conn, uri)
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
