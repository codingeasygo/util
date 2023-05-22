package xdebug

import (
	"io"
	"net"
)

// EchoServer
type EchoServer struct {
	net.Listener
	Address string
}

// NewEchoServer will return new EchoServer
func NewEchoServer(network, address string) (server *EchoServer, err error) {
	server = &EchoServer{}
	server.Listener, err = net.Listen(network, address)
	if err == nil {
		server.Address = server.Listener.Addr().String()
		go server.procAccept()
	}
	return
}

func (e *EchoServer) procAccept() {
	for {
		conn, err := e.Listener.Accept()
		if err != nil {
			break
		}
		go io.Copy(conn, conn)
	}
}
