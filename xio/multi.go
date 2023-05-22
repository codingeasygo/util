package xio

import "io"

// MultiWriteCloser is WriterCloser to bind one write to mulit sub writer
type MultiWriteCloser struct {
	Writers []io.Writer
}

// NewMultiWriter will return new MultiWriteCloser
func NewMultiWriter(writers ...io.Writer) (writer *MultiWriteCloser) {
	writer = &MultiWriteCloser{Writers: writers}
	return
}

func (m *MultiWriteCloser) Write(p []byte) (n int, err error) {
	for _, w := range m.Writers {
		n, err = w.Write(p)
		if err != nil {
			break
		}
		if n != len(p) {
			err = io.ErrShortWrite
			break
		}
	}
	return
}

// Close will try  close all writer if it is io.Closer
func (m *MultiWriteCloser) Close() (err error) {
	for _, w := range m.Writers {
		if closer, ok := w.(io.Closer); ok {
			xerr := closer.Close()
			if err == nil {
				err = xerr
			}
		}
	}
	return
}
