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
	if buffer.String() != "abc" {
		t.Error("error")
		return
	}
	//large write
	fmt.Fprintf(buffer, "xabc")
	if buffer.String() != "abc" {
		t.Error("error")
		return
	}
}
