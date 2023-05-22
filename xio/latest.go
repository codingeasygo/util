package xio

// LatestBuffer is membory buffer to remain latest buffer size bytes
type LatestBuffer struct {
	buffer []byte
	length int
}

// NewLatestBuffer will create remain bufffer
func NewLatestBuffer(bufferSize int) (buffer *LatestBuffer) {
	buffer = &LatestBuffer{
		buffer: make([]byte, bufferSize),
		length: 0,
	}
	return
}

// Bytes will get the having data
func (r *LatestBuffer) Bytes() []byte {
	return r.buffer[0:r.length]
}

func (r *LatestBuffer) Write(p []byte) (n int, err error) {
	if len(p) > len(r.buffer) {
		p = p[len(p)-len(r.buffer):]
	}
	offset := r.length
	if len(p)+r.length > len(r.buffer) {
		offset = len(r.buffer) - len(p)
	}
	if offset > 0 {
		copy(r.buffer[0:], r.buffer[r.length-offset:])
	}
	n = copy(r.buffer[offset:], p)
	r.length = offset + n
	return
}

func (r *LatestBuffer) String() string {
	return string(r.buffer[0:r.length])
}
