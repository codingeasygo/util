package xio

import (
	"fmt"
	"testing"
)

func TestLatestBuffer(t *testing.T) {
	buffer := NewLatestBuffer(3)
	//multi write
	fmt.Fprintf(buffer, "xa")
	fmt.Fprintf(buffer, "bc")
	if "abc" != string(buffer.Bytes()) {
		t.Error("error")
		return
	}
	//large write
	fmt.Fprintf(buffer, "xabc")
	if "abc" != string(buffer.Bytes()) {
		t.Error("error")
		return
	}
}
