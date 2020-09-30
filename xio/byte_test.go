package xio

import (
	"bytes"
	"fmt"
	"testing"
)

func TestByteDistributeWriteCloser(t *testing.T) {
	dist := NewByteDistributeWriter()
	buffer1 := bytes.NewBuffer(nil)
	buffer2 := bytes.NewBuffer(nil)
	dist.Add('a', NewCombinedReadWriteCloser(nil, buffer1, nil))
	dist.Add('b', NewCombinedReadWriteCloser(nil, buffer2, nil))
	fmt.Fprintf(dist, "a123")
	fmt.Fprintf(dist, "b000")
	if string(buffer1.Bytes()) != "a123" {
		t.Error("error")
		return
	}
	if string(buffer2.Bytes()) != "b000" {
		t.Error("error")
		return
	}
	_, err := dist.Write([]byte("111"))
	if err == nil {
		t.Error(err)
		return
	}
	dist.Close()
}
