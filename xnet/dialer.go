package xnet

import (
	"io"
	"net"
	"net/url"
)

// Dialer is interface for dial raw connect by string
type Dialer interface {
	Dial(remote string) (raw io.ReadWriteCloser, err error)
}

// RawDialer is an interface to dial raw conenction
type RawDialer interface {
	Dial(network, address string) (net.Conn, error)
}

// RawDialerF is an the implementation of RawDialer by func
type RawDialerF func(network, address string) (net.Conn, error)

// Dial dial to remote by func
func (d RawDialerF) Dial(network, address string) (raw net.Conn, err error) {
	raw, err = d(network, address)
	return
}

type RawDialerWrapper struct {
	RawDialer
}

func NewRawDialerWrapper(raw RawDialer) (dialer *RawDialerWrapper) {
	dialer = &RawDialerWrapper{
		RawDialer: raw,
	}
	return
}

func NewNetDailer() (dialer *RawDialerWrapper) {
	dialer = NewRawDialerWrapper(&net.Dialer{})
	return
}

func (w RawDialerWrapper) Dial(remote string) (raw io.ReadWriteCloser, err error) {
	remoteURI, err := url.Parse(remote)
	if err != nil {
		return
	}
	network := remoteURI.Scheme
	address := remoteURI.Host
	if len(remoteURI.Port()) < 1 {
		switch network {
		case "https":
			address += ":443"
			network = "tcp"
		case "http":
			address += ":80"
			network = "tcp"
		}
	}
	raw, err = w.RawDialer.Dial(network, address)
	return
}
