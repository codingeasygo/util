package xio

import (
	"fmt"
	"io"
)

// CombinedReadWriteCloser is middle struct to combine Reader/Writer/Closer
type CombinedReadWriteCloser struct {
	io.Reader
	io.Writer
	io.Closer
}

// NewCombinedReadWriteCloser will return new combined
func NewCombinedReadWriteCloser(reader io.Reader, writer io.Writer, closer io.Closer) (combined *CombinedReadWriteCloser) {
	combined = &CombinedReadWriteCloser{
		Reader: reader,
		Writer: writer,
		Closer: closer,
	}
	return
}

func (c *CombinedReadWriteCloser) Read(p []byte) (n int, err error) {
	if c.Reader == nil {
		err = fmt.Errorf("combined reader is nil")
		return
	}
	n, err = c.Reader.Read(p)
	return
}

func (c *CombinedReadWriteCloser) Write(p []byte) (n int, err error) {
	if c.Writer == nil {
		err = fmt.Errorf("combined writer is nil")
		return
	}
	n, err = c.Writer.Write(p)
	return
}

// Close will close Closer
func (c *CombinedReadWriteCloser) Close() (err error) {
	if c.Closer != nil {
		err = c.Closer.Close()
	}
	return
}
