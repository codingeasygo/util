package frame

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"
)

const (
	//OffsetBytes is default buffer offset
	OffsetBytes = 4
)

//ErrFrameTooLarge is the error when the frame head lenght > buffer length
var ErrFrameTooLarge = fmt.Errorf("%v", "frame is too large")

type readDeadlinable interface {
	SetReadDeadline(t time.Time) error
}

type writeDeadlinable interface {
	SetWriteDeadline(t time.Time) error
}

//Reader is interface for read the raw io as frame mode
type Reader interface {
	io.Reader
	ReadFrame() (frame []byte, err error)
	SetReadTimeout(timeout time.Duration)
}

//Writer is interface for write the raw io as frame mode
type Writer interface {
	io.Writer
	//WriteCmd will write data by frame mode, it must have 4 bytes at the begin of buffer to store the frame length.
	//genral buffer is (4 bytes)+(user data), 4 bytes will be set the in WriteCmd
	WriteFrame(buffer []byte) (n int, err error)
	SetWriteTimeout(timeout time.Duration)
}

//ReadWriter is interface for read/write the raw io as frame mode
type ReadWriter interface {
	Reader
	Writer
	SetTimeout(timeout time.Duration)
}

//ReadWriteCloser is interface for read/write the raw io as frame mode
type ReadWriteCloser interface {
	ReadWriter
	io.Closer
}

//NewReader will create new Reader by raw reader and buffer size
func NewReader(raw io.Reader, bufferSize int) (reader *BaseReader) {
	reader = NewBaseReader(raw, bufferSize)
	return
}

//NewWriter will return new BaseWriter
func NewWriter(raw io.Writer) (writer *BaseWriter) {
	writer = NewBaseWriter(raw)
	return
}

//BaseReadWriteCloser is frame reader/writer combiner
type BaseReadWriteCloser struct {
	io.Closer
	*BaseReader
	*BaseWriter
}

//Close will call the closer
func (b *BaseReadWriteCloser) Close() (err error) {
	if b.Closer != nil {
		err = b.Closer.Close()
	}
	return
}

func (b *BaseReadWriteCloser) String() string {
	var rawReader, rawWriter interface{} = b.BaseReader.Raw, b.BaseWriter.Raw
	if rawReader == rawWriter {
		return fmt.Sprintf("%v", b.BaseReader)
	}
	return fmt.Sprintf("Reader:%v,Writer:%v", b.BaseReader, b.BaseWriter)
}

//SetTimeout will record the timout
func (b *BaseReadWriteCloser) SetTimeout(timeout time.Duration) {
	b.BaseReader.SetReadTimeout(timeout)
	b.BaseWriter.SetWriteTimeout(timeout)
}

//NewReadWriter will return new ReadWriteCloser
func NewReadWriter(raw io.ReadWriter, bufferSize int) (frame *BaseReadWriteCloser) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	closer, _ := raw.(io.Closer)
	frame = &BaseReadWriteCloser{
		Closer:     closer,
		BaseReader: NewBaseReader(raw, bufferSize),
		BaseWriter: NewBaseWriter(raw),
	}
	return
}

//NewReadWriteCloser will return new ReadWriteCloser
func NewReadWriteCloser(raw io.ReadWriteCloser, bufferSize int) (frame *BaseReadWriteCloser) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	frame = &BaseReadWriteCloser{
		Closer:     raw,
		BaseReader: NewBaseReader(raw, bufferSize),
		BaseWriter: NewBaseWriter(raw),
	}
	return
}

//BaseReader imple read raw connection by frame mode
type BaseReader struct {
	Buffer  []byte
	Raw     io.Reader
	Timeout time.Duration
	offset  uint32
	length  uint32
	locker  sync.RWMutex
}

//NewBaseReader will create new Reader by raw reader and buffer size
func NewBaseReader(raw io.Reader, bufferSize int) (reader *BaseReader) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	reader = &BaseReader{
		Buffer: make([]byte, bufferSize),
		Raw:    raw,
		locker: sync.RWMutex{},
	}
	return
}

//readMore will read more data to buffer
func (b *BaseReader) readMore() (err error) {
	if r, ok := b.Raw.(readDeadlinable); b.Timeout > 0 && ok {
		r.SetReadDeadline(time.Now().Add(b.Timeout))
	}
	readed, err := b.Raw.Read(b.Buffer[b.offset+b.length:])
	if err == nil {
		b.length += uint32(readed)
	}
	return
}

//ReadFrame will read raw reader as frame mode. it will return length(4bytes)+data.
//the return []byte is the buffer slice, must be copy to new []byte, it will be change after next read
func (b *BaseReader) ReadFrame() (cmd []byte, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	more := b.length < 5
	for {
		if more {
			err = b.readMore()
			if err != nil {
				break
			}
			if b.length < 5 {
				continue
			}
		}
		b.Buffer[b.offset] = 0
		frameLength := binary.BigEndian.Uint32(b.Buffer[b.offset:])
		if frameLength > uint32(len(b.Buffer)) {
			err = ErrFrameTooLarge
			break
		}
		if b.length < frameLength {
			more = true
			if b.offset > 0 {
				copy(b.Buffer[0:], b.Buffer[b.offset:b.offset+b.length])
				b.offset = 0
			}
			continue
		}
		cmd = b.Buffer[b.offset : b.offset+frameLength]
		cmd[0] = 0
		b.offset += frameLength
		b.length -= frameLength
		more = b.length <= 4
		if b.length < 1 {
			b.offset = 0
		}
		if more && b.offset > 0 {
			copy(b.Buffer[0:], b.Buffer[b.offset:b.offset+b.length])
			b.offset = 0
		}
		break
	}
	return
}

//Read implment the io.Reader
//it will read the one frame and copy the data to p
func (b *BaseReader) Read(p []byte) (n int, err error) {
	data, err := b.ReadFrame()
	if err == nil {
		n = copy(p, data[4:])
	}
	return
}

//SetReadTimeout will record the timout
func (b *BaseReader) SetReadTimeout(timeout time.Duration) {
	b.Timeout = timeout
}

func (b *BaseReader) String() string {
	return fmt.Sprintf("%v", b.Raw)
}

//BaseWriter implment the frame Writer
type BaseWriter struct {
	//the raw io writer
	Raw     io.Writer
	Timeout time.Duration
	locker  sync.RWMutex
}

//NewBaseWriter will return new BaseWriter
func NewBaseWriter(raw io.Writer) (writer *BaseWriter) {
	writer = &BaseWriter{Raw: raw}
	return
}

//WriteFrame will write data by frame mode, it must have 4 bytes at the begin of buffer to store the frame length.
//genral buffer is (4 bytes)+(user data), 4 bytes will be set the in WriteCmd
func (b *BaseWriter) WriteFrame(buffer []byte) (w int, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	binary.BigEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[0] = byte(rand.Intn(255))
	if w, ok := b.Raw.(writeDeadlinable); b.Timeout > 0 && ok {
		w.SetWriteDeadline(time.Now().Add(b.Timeout))
	}
	w, err = b.Raw.Write(buffer)
	return
}

//Write implment the io.Writer, the p is user data buffer.
//it will make a new []byte with len(p)+4, the copy data to buffer
func (b *BaseWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p)+4)
	copy(buf[4:], p)
	n = len(buf)
	_, err = b.WriteFrame(buf)
	return
}

//SetWriteTimeout will record the timout
func (b *BaseWriter) SetWriteTimeout(timeout time.Duration) {
	b.Timeout = timeout
}

func (b *BaseWriter) String() string {
	return fmt.Sprintf("%v", b.Raw)
}

//Wrap will create buffer and add frame header
func Wrap(p []byte) (buffer []byte) {
	buffer = make([]byte, len(p)+4)
	copy(buffer[4:], p)
	binary.BigEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[0] = byte(rand.Intn(255))
	return
}
