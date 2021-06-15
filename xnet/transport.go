package xnet

import (
	"fmt"
	"io"
	"strings"
)

type Transporter interface {
	Transport(conn io.ReadWriteCloser, remote string) (err error)
}

func transportCopy(a, b io.ReadWriteCloser) (aerr, berr error) {
	go func() {
		_, berr = io.Copy(a, b)
		a.Close()
	}()
	_, aerr = io.Copy(b, a)
	a.Close()
	return
}

func (w *WebsocketDialer) Transport(conn io.ReadWriteCloser, remote string) (err error) {
	raw, err := w.Dial(remote)
	if err != nil {
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
	parts := strings.SplitN(remote, "://", 2)
	if len(parts) < 2 {
		err = fmt.Errorf("invalid remote %v", remote)
		return
	}
	raw, err := d.Dial(parts[0], parts[1])
	if err != nil {
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
