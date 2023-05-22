package xio

import (
	"io"
)

type DiscardReadWriteCloser struct {
}

func NewDiscardReadWriteCloser() (discard *DiscardReadWriteCloser) {
	discard = &DiscardReadWriteCloser{}
	return
}

func (d *DiscardReadWriteCloser) Read(p []byte) (n int, err error) {
	err = io.EOF
	return
}

func (c *DiscardReadWriteCloser) Write(p []byte) (n int, err error) {
	n = len(p)
	return
}

func (c *DiscardReadWriteCloser) Close() (err error) {
	return
}
