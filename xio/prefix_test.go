package xio

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrefixReader(t *testing.T) {
	buffer := make([]byte, 1024)
	reader := NewPrefixReader(bytes.NewBufferString("abc"))
	reader.PreRead(1)
	err := FullBuffer(reader, buffer, 3, nil)
	if err != nil || string(buffer[0:3]) != "abc" {
		t.Error(err)
		return
	}
}

func TestPrefixReadWriteCloser(t *testing.T) {
	conn := NewPrefixReadWriteCloser(NewEchoConn())
	fmt.Fprintf(conn, "abc")
	conn.PreRead(1)
	buffer := make([]byte, 1024)
	err := FullBuffer(conn, buffer, 3, nil)
	if err != nil || string(buffer[0:3]) != "abc" {
		t.Error(err)
		return
	}
}

func TestNewPrefixConn(t *testing.T) {
	conn := NewPrefixConn(NewEchoConn())
	fmt.Fprintf(conn, "abc")
	conn.PreRead(1)
	buffer := make([]byte, 1024)
	err := FullBuffer(conn, buffer, 3, nil)
	if err != nil || string(buffer[0:3]) != "abc" {
		t.Error(err)
		return
	}
}
