package proxy

import (
	"net"

	"github.com/codingeasygo/util/proxy/http"
	"github.com/codingeasygo/util/proxy/socks"
	"github.com/codingeasygo/util/xio"
)

//Server provider http/socks combined server
type Server struct {
	*xio.ByteDistributeProcessor
	Dialer xio.PiperDialer
	HTTP   *http.Server
	SOCKS  *socks.Server
}

//NewServer will return new Server
func NewServer(dialer xio.PiperDialer) (server *Server) {
	server = &Server{
		ByteDistributeProcessor: xio.NewByteDistributeProcessor(),
		Dialer:                  dialer,
		HTTP:                    http.NewServer(),
		SOCKS:                   socks.NewServer(),
	}
	server.HTTP.Dialer = server
	server.SOCKS.Dialer = server
	server.AddProcessor('H', server.HTTP)
	server.AddProcessor('*', server.SOCKS)
	return
}

//DialPiper is xio.Piper implement
func (s *Server) DialPiper(uri string, bufferSize int) (raw xio.Piper, err error) {
	raw, err = s.Dialer.DialPiper(uri, bufferSize)
	return
}

//Start wiil listen tcp on addr and run process accept to ByteDistributeProcessor
func (s *Server) Start(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	go s.ProcAccept(listener)
	return
}
