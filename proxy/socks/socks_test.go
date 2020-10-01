package socks

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/codingeasygo/util/xio"
)

func TestProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "abc")
	}))
	server := NewServer()
	err := server.Start(":8011")
	if err != nil {
		t.Error(err)
		return
	}
	{ //CONNECT
		client := http.Client{
			Transport: &http.Transport{
				Dial: func(network, address string) (conn net.Conn, err error) {
					conn, err = Dial(":8011", address)
					return
				},
			},
		}
		{ //ok
			resp, err := client.Get(ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if string(data) != "abc" {
				t.Error("error")
				return
			}
		}
		{ //error
			_, err := client.Get("http://127.0.0.1:233")
			if err == nil {
				t.Error(err)
				return
			}
		}
	}
	{ //ERROR
		conn, _ := net.Dial("tcp", ":8011")
		time.Sleep(10 * time.Millisecond)
		conn.Close()
		Dial("127.0.0.1:233", "")
	}
	server.Stop()
	go func() {
		server.Run(":8011")
	}()
	time.Sleep(10 * time.Millisecond)
	server.Stop()
}

func proxyDial(t *testing.T, remote string, port uint16) {
	conn, err := net.Dial("tcp", "localhost:2081")
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	buf := make([]byte, 1024*64)
	proxyReader := bufio.NewReader(conn)
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return
	}
	err = xio.FullBuffer(proxyReader, buf, 2, nil)
	if err != nil {
		return
	}
	if buf[0] != 0x05 || buf[1] != 0x00 {
		err = fmt.Errorf("only ver 0x05 / method 0x00 is supported, but %x/%x", buf[0], buf[1])
		return
	}
	buf[0], buf[1], buf[2], buf[3] = 0x05, 0x01, 0x00, 0x03
	buf[4] = byte(len(remote))
	copy(buf[5:], []byte(remote))
	binary.BigEndian.PutUint16(buf[5+len(remote):], port)
	_, err = conn.Write(buf[:buf[4]+7])
	if err != nil {
		return
	}
	readed, err := proxyReader.Read(buf)
	if err != nil {
		return
	}
	fmt.Printf("->%v\n", buf[0:readed])
	fmt.Fprintf(conn, "abc")
}

func proxyDial2(t *testing.T, remote string, port uint16) {
	conn, err := net.Dial("tcp", "localhost:2081")
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	buf := make([]byte, 1024*64)
	proxyReader := bufio.NewReader(conn)
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return
	}
	err = xio.FullBuffer(proxyReader, buf, 2, nil)
	if err != nil {
		return
	}
	if buf[0] != 0x05 || buf[1] != 0x00 {
		err = fmt.Errorf("only ver 0x05 / method 0x00 is supported, but %x/%x", buf[0], buf[1])
		return
	}
	buf[0], buf[1], buf[2], buf[3] = 0x05, 0x01, 0x00, 0x13
	buf[4] = byte(len(remote))
	copy(buf[5:], []byte(remote))
	binary.BigEndian.PutUint16(buf[5+len(remote):], port)
	_, err = conn.Write(buf[:buf[4]+7])
	if err != nil {
		return
	}
	readed, err := proxyReader.Read(buf)
	if err != nil {
		return
	}
	fmt.Printf("->%v\n", buf[0:readed])
	fmt.Fprintf(conn, "abc")
}

func proxyDialIP(t *testing.T, bys []byte, port uint16) {
	conn, err := net.Dial("tcp", "localhost:2081")
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	buf := make([]byte, 1024*64)
	proxyReader := bufio.NewReader(conn)
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return
	}
	err = xio.FullBuffer(proxyReader, buf, 2, nil)
	if err != nil {
		return
	}
	if buf[0] != 0x05 || buf[1] != 0x00 {
		err = fmt.Errorf("only ver 0x05 / method 0x00 is supported, but %x/%x", buf[0], buf[1])
		return
	}
	buf[0], buf[1], buf[2], buf[3] = 0x05, 0x01, 0x00, 0x01
	copy(buf[4:], bys)
	binary.BigEndian.PutUint16(buf[8:], port)
	_, err = conn.Write(buf[:10])
	if err != nil {
		return
	}
	readed, err := proxyReader.Read(buf)
	if err != nil {
		return
	}
	fmt.Printf("->%v\n", buf[0:readed])
}

func proxyDialIPv6(t *testing.T, bys []byte, port uint16) {
	conn, err := net.Dial("tcp", "localhost:2081")
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	buf := make([]byte, 1024*64)
	proxyReader := bufio.NewReader(conn)
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return
	}
	err = xio.FullBuffer(proxyReader, buf, 2, nil)
	if err != nil {
		return
	}
	if buf[0] != 0x05 || buf[1] != 0x00 {
		err = fmt.Errorf("only ver 0x05 / method 0x00 is supported, but %x/%x", buf[0], buf[1])
		return
	}
	buf[0], buf[1], buf[2], buf[3] = 0x05, 0x01, 0x00, 0x04
	copy(buf[4:], bys)
	binary.BigEndian.PutUint16(buf[8:], port)
	_, err = conn.Write(buf[:10])
	if err != nil {
		return
	}
	readed, err := proxyReader.Read(buf)
	if err != nil {
		return
	}
	fmt.Printf("->%v\n", buf[0:readed])
}

type CodableErr struct {
	Err error
}

func (c *CodableErr) Error() string {
	return c.Err.Error()
}

func (c *CodableErr) Code() byte {
	return 0x01
}

func TestSocksProxy(t *testing.T) {
	proxy := NewServer()
	proxy.Dialer = func(utype int, uri string, conn io.ReadWriteCloser) (raw io.ReadWriteCloser, err error) {
		r, w, _ := os.Pipe()
		raw = xio.NewCombinedReadWriteCloser(r, w, w)
		// raw, err = net.Dial("tcp", uri)
		// if err == nil {
		// 	go io.Copy(conn, raw)
		// 	go io.Copy(raw, conn)
		// 	time.Sleep(100 * time.Millisecond)
		// } else {
		// 	err = &CodableErr{Err: err}
		// }
		// fmt.Println("dial to ", uri, err)
		return
	}
	go func() {
		proxy.Run(":2081")
	}()
	proxyDial(t, "localhost", 80)
	proxyDial2(t, "localhost:80", 0)
	proxyDial(t, "localhost", 81)
	proxyDialIP(t, make([]byte, 4), 80)
	// proxyDialIPv6(t, make([]byte, 16), 80)
	{ //test error
		//
		conn, conb, _ := xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Close()
		//
		conn, conb, _ = xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Write([]byte{0x00, 0x00})
		conn.Close()
		//
		conn, conb, _ = xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Write([]byte{0x05, 0x01})
		conn.Close()
		//
		conn, conb, _ = xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Write([]byte{0x05, 0x01, 0x00})
		conn.Close()
		//
		conn, conb, _ = xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Write([]byte{0x05, 0x01, 0x00})
		conn.Read(make([]byte, 1024))
		conn.Close()
		//
		conn, conb, _ = xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Write([]byte{0x05, 0x01, 0x00})
		conn.Read(make([]byte, 1024))
		conn.Write([]byte{0x00, 0x01, 0x00, 0x00, 0x00})
		conn.Close()
		//
		conn, conb, _ = xio.CreatePipedConn()
		go proxy.ProcConn(conb)
		conn.Write([]byte{0x05, 0x01, 0x00})
		conn.Read(make([]byte, 1024))
		buf := []byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x010}
		binary.BigEndian.PutUint16(buf[8:], 80)
		conn.Write(buf)
		conn.Close()
		time.Sleep(time.Second)
	}
	proxy.Stop()
}
