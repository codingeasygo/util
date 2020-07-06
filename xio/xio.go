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

//CopyMulti will copy data from Reader and write to multi Writer at the same time
func CopyMulti(dst []io.Writer, src io.Reader) (written int64, err error) {
	written, err = CopyBufferMulti(dst, src, nil)
	return
}

//CopyBufferMulti will copy data from Reader and write to multi Writer at the same time
func CopyBufferMulti(dst []io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf == nil {
		size := 32 * 1024
		buf = make([]byte, size)
	}
	write := func(nr int, b []byte) (nw int, err error) {
		for _, d := range dst {
			nw, err = d.Write(b)
			if err != nil {
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		return
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := write(nr, buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

//CopyMax will copy data to writer and total limit by max
func CopyMax(dst io.Writer, src io.Reader, max int64) (written int64, err error) {
	written, err = CopyBufferMax(dst, src, max, nil)
	return
}

//CopyBufferMax will copy data to writer and total limit by max
func CopyBufferMax(dst io.Writer, src io.Reader, max int64, buf []byte) (written int64, err error) {
	if buf == nil {
		size := 32 * 1024
		buf = make([]byte, size)
	}
	for {
		limited := max - written
		if limited < 1 {
			err = fmt.Errorf("copy max limit")
			break
		}
		if limited > int64(len(buf)) {
			limited = int64(len(buf))
		}
		nr, er := src.Read(buf[0:limited])
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
