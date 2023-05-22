package xio

import (
	"fmt"
	"io"
	"sync"
)

// ByteDistributeWriter is Writer by byte
type ByteDistributeWriter struct {
	ws  map[byte]io.Writer
	lck sync.RWMutex
}

// NewByteDistributeWriter will return new ByteDistributeWriter
func NewByteDistributeWriter() (writer *ByteDistributeWriter) {
	writer = &ByteDistributeWriter{
		ws:  map[byte]io.Writer{},
		lck: sync.RWMutex{},
	}
	return
}

// Add WriteCloser to list
func (h *ByteDistributeWriter) Add(m byte, w io.Writer) {
	h.lck.Lock()
	defer h.lck.Unlock()
	h.ws[m] = w
}

func (h *ByteDistributeWriter) Write(b []byte) (n int, err error) {
	h.lck.RLock()
	writer, ok := h.ws[b[0]]
	h.lck.RUnlock()
	if ok {
		n, err = writer.Write(b)
	} else {
		err = fmt.Errorf("writer not exist by %v", b[0])
	}
	return
}

// Close will close all connection
func (h *ByteDistributeWriter) Close() (err error) {
	h.lck.Lock()
	for k, w := range h.ws {
		if closer, ok := w.(io.Closer); ok {
			closer.Close()
		}
		delete(h.ws, k)
	}
	h.lck.Unlock()
	return
}
