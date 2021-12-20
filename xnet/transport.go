package xnet

import (
	"io"
	"net/url"
)

type Transporter interface {
	Transport(conn io.ReadWriteCloser, remote string) (err error)
}

func transportCopy(a, b io.ReadWriteCloser) (aerr, berr error) {
	go func() {
		_, berr = io.Copy(a, b)
		a.Close()
		b.Close()
	}()
	_, aerr = io.Copy(b, a)
	a.Close()
	b.Close()
	return
}

func (w *WebsocketDialer) Transport(conn io.ReadWriteCloser, remote string) (err error) {
	raw, err := w.Dial(remote)
	if err != nil {
		conn.Close()
		return
	}
	aerr, berr := transportCopy(conn, raw)
	if aerr != nil {
		err = aerr
	} else {
		err = berr
	}
	return
}

func (d RawDialerF) Transport(conn io.ReadWriteCloser, remote string) (err error) {
	u, err := url.Parse(remote)
	if err != nil {
		conn.Close()
		return
	}
	raw, err := d.Dial(u.Scheme, u.Host)
	if err != nil {
		conn.Close()
		return
	}
	aerr, berr := transportCopy(conn, raw)
	if aerr != nil {
		err = aerr
	} else {
		err = berr
	}
	return
}
