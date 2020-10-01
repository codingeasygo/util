package proxy

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"

	"github.com/codingeasygo/util/xio"
)

func init() {
	go http.ListenAndServe(":6060", nil)
}

func TestByteDistribute(t *testing.T) {
	accept := make(chan net.Conn, 1)
	processor := NewByteDistributeProcessor()
	go processor.ProcAccept(xio.ListenerF(func() (conn net.Conn, err error) {
		conn = <-accept
		if conn == nil {
			err = fmt.Errorf("closed")
		}
		return
	}))
	processor.AddProcessor('A', ProcessorF(func(conn net.Conn) (err error) {
		buf := make([]byte, 1024)
		for {
			readed, err := conn.Read(buf)
			if err != nil {
				break
			}

			fmt.Fprintf(conn, "A:%v", string(buf[1:readed]))
		}
		return
	}))
	processor.AddProcessor('*', ProcessorF(func(conn net.Conn) (err error) {
		buf := make([]byte, 1024)
		for {
			readed, err := conn.Read(buf)
			if err != nil {
				break
			}
			fmt.Fprintf(conn, "*:%v", string(buf[1:readed]))
		}
		return
	}))
	buf := make([]byte, 1024)
	//
	//A
	conna, connb, _ := xio.CreatePipedConn()
	accept <- connb
	fmt.Fprintf(conna, "A123")
	readed, err := conna.Read(buf)
	if err != nil || string(buf[0:readed]) != "A:123" {
		t.Error(err)
		return
	}
	conna.Close()
	//
	//B
	conna, connb, _ = xio.CreatePipedConn()
	accept <- connb
	fmt.Fprintf(conna, "B123")
	readed, err = conna.Read(buf)
	if err != nil || string(buf[0:readed]) != "*:123" {
		t.Error(err)
		return
	}
	// conna.Close()
	//
	//error
	conna, connb, _ = xio.CreatePipedConn()
	accept <- connb
	conna.Close()
	//
	//
	processor.Close()
	accept <- nil
	time.Sleep(10 * time.Millisecond)
}
