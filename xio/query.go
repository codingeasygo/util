package xio

import (
	"context"
	"fmt"
	"net"
	"time"
)

type QueryConn struct {
	sendData chan []byte
	recvData []byte
	recvWait chan int
	err      error
}

func NewQueryConn() (conn *QueryConn) {
	conn = &QueryConn{
		sendData: make(chan []byte, 1),
		recvWait: make(chan int, 1),
	}
	return
}

func (q *QueryConn) Read(p []byte) (n int, err error) {
	data := <-q.sendData
	if len(data) < 1 {
		err = q.err
		return
	}
	n = copy(p, data)
	return
}

func (q *QueryConn) Write(p []byte) (n int, err error) {
	if len(q.recvData) > 0 {
		err = fmt.Errorf("recv twice")
		return
	}
	q.recvData = make([]byte, len(p))
	n = copy(q.recvData, p)
	select {
	case q.recvWait <- 1:
	default:
	}
	return
}

func (q *QueryConn) Close() (err error) {
	q.err = fmt.Errorf("closed")
	select {
	case q.sendData <- nil:
	default:
	}
	select {
	case q.recvWait <- 1:
	default:
	}
	return
}

func (q *QueryConn) clearSend() {
	select {
	case <-q.sendData:
	default:
	}
}

func (q *QueryConn) Query(ctx context.Context, request []byte) (response []byte, err error) {
	select {
	case q.sendData <- request:
		err = q.err
	case <-ctx.Done():
		err = fmt.Errorf("context canceled")
	}
	if err != nil {
		q.clearSend()
		return
	}
	select {
	case <-q.recvWait:
		err = q.err
		response = q.recvData
		q.recvData = nil
	case <-ctx.Done():
		err = fmt.Errorf("context canceled")
	}
	return
}

// LocalAddr returns the local network address.
func (q *QueryConn) LocalAddr() net.Addr {
	return q
}

// RemoteAddr returns the remote network address.
func (q *QueryConn) RemoteAddr() net.Addr {
	return q
}

// SetDeadline for net.Conn
func (q *QueryConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline for net.Conn
func (q *QueryConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline for net.Conn
func (q *QueryConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// Network is net.Addr impl
func (q *QueryConn) Network() string {
	return "query"
}

func (q *QueryConn) String() string {
	return "query"
}
