package xio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	WriteJSON(ioutil.Discard, map[string]interface{}{})
}

func TestCopyPacketConn(t *testing.T) {
	packet1, _ := net.ListenPacket("udp", ":12332")
	go CopyPacketConn(packet1, packet1)
	raw, err := net.Dial("udp", "127.0.0.1:12332")
	if err != nil {
		t.Error(err)
		return
	}
	defer raw.Close()
	udp := raw.(*net.UDPConn)
	_, err = udp.Write([]byte("123"))
	if err != nil {
		t.Error(err)
		return
	}
	buffer := make([]byte, 123)
	n, err := udp.Read(buffer)
	if err != nil {
		t.Error(err)
		return
	}
	if string(buffer[0:n]) != "123" {
		t.Error("error")
		return
	}
	packet1.Close()
	time.Sleep(100 * time.Millisecond)
}

func TestCopyPacketConnToWriter(t *testing.T) {
	packet1, _ := net.ListenPacket("udp", ":12332")
	buffer := bytes.NewBuffer(nil)
	go CopyPacketConn(buffer, packet1)
	raw, err := net.Dial("udp", "127.0.0.1:12332")
	if err != nil {
		t.Error(err)
		return
	}
	defer raw.Close()
	udp := raw.(*net.UDPConn)
	_, err = udp.Write([]byte("123"))
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	if string(buffer.Bytes()) != "123" {
		t.Error("error")
		return
	}
	packet1.Close()
	time.Sleep(100 * time.Millisecond)
}

func TestCopyPacketConnErrorNotSupported(t *testing.T) {
	packet1, _ := net.ListenPacket("udp", ":12332")
	raw, err := net.Dial("udp", "127.0.0.1:12332")
	if err != nil {
		t.Error(err)
		return
	}
	defer raw.Close()
	udp := raw.(*net.UDPConn)
	_, err = udp.Write([]byte("123"))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = CopyPacketConn(t, packet1)
	if err == nil {
		t.Error(nil)
	}
	packet1.Close()
	time.Sleep(100 * time.Millisecond)
}

func TestCopyPacketTo(t *testing.T) {
	packet1, _ := net.ListenPacket("udp", ":12332")
	raw, err := net.Dial("udp", "127.0.0.1:12332")
	if err != nil {
		t.Error(err)
		return
	}
	defer raw.Close()

	//
	udp := raw.(*net.UDPConn)
	_, err = CopyPacketTo(packet1, udp.LocalAddr(), bytes.NewBuffer([]byte("123")))
	if err != io.EOF {
		t.Error(err)
		return
	}
	buffer := make([]byte, 1024)
	n, err := udp.Read(buffer)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	if string(buffer[0:n]) != "123" {
		t.Error("error")
		return
	}
	packet1.Close()
	time.Sleep(100 * time.Millisecond)
}

func TestCopyPacketToError(t *testing.T) {
	packet1, _ := net.ListenPacket("udp", ":12332")
	packet1.Close()
	_, err := CopyPacketTo(packet1, packet1.LocalAddr(), bytes.NewBuffer([]byte("123")))
	if err == nil {
		t.Error(err)
		return
	}
}

type copyMultiTestReader struct {
}

func (c *copyMultiTestReader) Read(p []byte) (n int, err error) {
	err = fmt.Errorf("test error")
	return
}

type copyMultiTestWriter struct {
	n int
}

func (c *copyMultiTestWriter) Write(p []byte) (n int, err error) {
	switch c.n {
	case 0:
		n = len(p)
	case 1:
		n = len(p) - 1
	case 2:
		err = fmt.Errorf("test error")
	}
	return
}

func TestCopyMulti(t *testing.T) {
	var err error
	writer1 := &copyMultiTestWriter{}
	writer2 := &copyMultiTestWriter{}
	//
	writer1.n, writer2.n = 0, 0
	_, err = CopyMulti([]io.Writer{writer1, writer2}, bytes.NewBufferString("123"))
	if err != nil {
		t.Error(err)
		return
	}
	//
	writer1.n, writer2.n = 0, 0
	_, err = CopyMulti([]io.Writer{writer1, writer2}, &copyMultiTestReader{})
	if err == nil {
		t.Error(err)
		return
	}
	//
	writer1.n, writer2.n = 1, 1
	_, err = CopyMulti([]io.Writer{writer1, writer2}, bytes.NewBufferString("123"))
	if err != io.ErrShortWrite {
		t.Error(err)
		return
	}
	//
	writer1.n, writer2.n = 2, 2
	_, err = CopyMulti([]io.Writer{writer1, writer2}, bytes.NewBufferString("123"))
	if err == nil {
		t.Error(err)
		return
	}
}

func TestCopyMax(t *testing.T) {
	var err error
	//
	_, err = CopyMax(bytes.NewBuffer(nil), bytes.NewBufferString("abc"), 10)
	if err != nil {
		t.Error(err)
		return
	}
	//
	_, err = CopyMax(bytes.NewBuffer(nil), bytes.NewBufferString("abc"), 2)
	if err == nil {
		t.Error(err)
		return
	}
	//
	_, err = CopyBufferMax(bytes.NewBuffer(nil), bytes.NewBufferString("abc"), 10, make([]byte, 2))
	if err != nil {
		t.Error(err)
		return
	}
	//
	writer1 := &copyMultiTestWriter{}
	//
	writer1.n = 1
	_, err = CopyMax(writer1, bytes.NewBufferString("abc"), 1024)
	if err == nil {
		t.Error(err)
		return
	}
	//
	writer1.n = 2
	_, err = CopyMax(writer1, bytes.NewBufferString("abc"), 1024)
	if err == nil {
		t.Error(err)
		return
	}
	//
	writer1.n = 0
	_, err = CopyMax(ioutil.Discard, &copyMultiTestReader{}, 1024)
	if err == nil {
		t.Error(err)
		return
	}
}
