package xio

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
)

//WriteJSON will marshal value to string and write to io.Writer
func WriteJSON(w io.Writer, v interface{}) (n int, err error) {
	data, err := json.Marshal(v)
	if err == nil {
		n, err = w.Write(data)
	}
	return
}

func CopyPacketConn(dst interface{}, src net.PacketConn) (l int64, err error) {
	buffer := make([]byte, 2*1024)
	for {
		n, from, xerr := src.ReadFrom(buffer)
		if xerr != nil {
			err = xerr
			break
		}
		if out, ok := dst.(net.PacketConn); ok {
			n, xerr = out.WriteTo(buffer[0:n], from)
		} else if out, ok := dst.(io.Writer); ok {
			n, xerr = out.Write(buffer[0:n])
		} else {
			xerr = fmt.Errorf("not supported dst by type %v", reflect.TypeOf(dst))
		}
		if xerr != nil {
			err = xerr
			break
		}
		l += int64(n)
	}
	return
}

func CopyPacketTo(dst net.PacketConn, to net.Addr, src io.Reader) (l int64, err error) {
	buffer := make([]byte, 2*1024)
	for {
		n, xerr := src.Read(buffer)
		if xerr != nil {
			err = xerr
			break
		}
		n, xerr = dst.WriteTo(buffer[0:n], to)
		if xerr != nil {
			err = xerr
			break
		}
		l += int64(n)
	}
	return
}
