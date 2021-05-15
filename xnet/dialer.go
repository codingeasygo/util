package xnet

import "net"

//RawDialer is an interface to dial raw conenction
type RawDialer interface {
	Dial(network, address string) (net.Conn, error)
}

//RawDialerF is an the implementation of RawDialer by func
type RawDialerF func(network, address string) (net.Conn, error)

//Dial dial to remote by func
func (d RawDialerF) Dial(network, address string) (raw net.Conn, err error) {
	raw, err = d(network, address)
	return
}
