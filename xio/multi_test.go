package xio

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestMultiWriter(t *testing.T) {
	var n int
	var err error
	raw := &copyMultiTestWriter{}
	writer := NewMultiWriter(ioutil.Discard, raw)
	raw.n = 0
	n, err = fmt.Fprintf(writer, "abc")
	if err != nil || n != 3 {
		t.Error(err)
		return
	}
	raw.n = 1
	_, err = fmt.Fprintf(writer, "abc")
	if err == nil {
		t.Error(err)
		return
	}
	raw.n = 2
	_, err = fmt.Fprintf(writer, "abc")
	if err == nil {
		t.Error(err)
		return
	}
	writer.Close()
}
