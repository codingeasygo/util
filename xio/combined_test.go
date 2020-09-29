package xio

import "testing"

func TestCombinedReadWriteCloser(t *testing.T) {
	reader := &copyMultiTestReader{}
	writer := &copyMultiTestWriter{}
	combined := NewCombinedReadWriteCloser(reader, writer, writer)
	combined.Close()
	combined.Read(nil)
	combined.Write(nil)
	combined.Reader = nil
	combined.Writer = nil
	combined.Read(nil)
	combined.Write(nil)
}
