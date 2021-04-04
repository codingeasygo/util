package frame

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/codingeasygo/util/xio"
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
	GetReadByteOrder() (order binary.ByteOrder)
	GetReadLengthFieldMagic() (value int)
	GetReadLengthFieldOffset() (value int)
	GetReadLengthFieldLength() (value int)
	GetReadLengthAdjustment() (value int)
	SetReadByteOrder(order binary.ByteOrder)
	SetReadLengthFieldMagic(value int)
	SetReadLengthFieldOffset(value int)
	SetReadLengthFieldLength(value int)
	SetReadLengthAdjustment(value int)
}

//Writer is interface for write the raw io as frame mode
type Writer interface {
	io.Writer
	//WriteCmd will write data by frame mode, it must have 4 bytes at the begin of buffer to store the frame length.
	//genral buffer is (4 bytes)+(user data), 4 bytes will be set the in WriteCmd
	WriteFrame(buffer []byte) (n int, err error)
	SetWriteTimeout(timeout time.Duration)
	GetWriteByteOrder() (order binary.ByteOrder)
	GetWriteLengthFieldMagic() (value int)
	GetWriteLengthFieldOffset() (value int)
	GetWriteLengthFieldLength() (value int)
	GetWriteLengthAdjustment() (value int)
	SetWriteByteOrder(order binary.ByteOrder)
	SetWriteLengthFieldMagic(value int)
	SetWriteLengthFieldOffset(value int)
	SetWriteLengthFieldLength(value int)
	SetWriteLengthAdjustment(value int)
}

//ReadWriter is interface for read/write the raw io as frame mode
type ReadWriter interface {
	Reader
	Writer
	SetTimeout(timeout time.Duration)
	GetByteOrder() (order binary.ByteOrder)
	GetLengthFieldMagic() (value int)
	GetLengthFieldOffset() (value int)
	GetLengthFieldLength() (value int)
	GetLengthAdjustment() (value int)
	SetByteOrder(order binary.ByteOrder)
	SetLengthFieldMagic(value int)
	SetLengthFieldOffset(value int)
	SetLengthFieldLength(value int)
	SetLengthAdjustment(value int)
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
	return fmt.Sprintf("Reader:%v,Writer:%v", b.BaseReader, b.BaseWriter)
}

//SetTimeout will record the timout
func (b *BaseReadWriteCloser) SetTimeout(timeout time.Duration) {
	b.BaseReader.SetReadTimeout(timeout)
	b.BaseWriter.SetWriteTimeout(timeout)
}

func (b *BaseReadWriteCloser) GetByteOrder() (order binary.ByteOrder) {
	order = b.BaseReader.GetReadByteOrder()
	return
}

func (b *BaseReadWriteCloser) GetLengthFieldMagic() (value int) {
	value = b.BaseReader.GetReadLengthFieldMagic()
	return
}

func (b *BaseReadWriteCloser) GetLengthFieldOffset() (value int) {
	value = b.BaseReader.GetReadLengthFieldOffset()
	return
}

func (b *BaseReadWriteCloser) GetLengthFieldLength() (value int) {
	value = b.BaseReader.GetReadLengthFieldLength()
	return
}

func (b *BaseReadWriteCloser) GetLengthAdjustment() (value int) {
	value = b.BaseReader.GetReadLengthAdjustment()
	return
}

func (b *BaseReadWriteCloser) SetByteOrder(order binary.ByteOrder) {
	b.BaseReader.SetReadByteOrder(order)
	b.BaseWriter.SetWriteByteOrder(order)
}

//SetLengthFieldMagic will set the LengthFieldMagic for reader/writer
func (b *BaseReadWriteCloser) SetLengthFieldMagic(value int) {
	b.BaseReader.SetReadLengthFieldMagic(value)
	b.BaseWriter.SetWriteLengthFieldMagic(value)
}

//SetLengthFieldOffset will set the LengthFieldOffset for reader/writer
func (b *BaseReadWriteCloser) SetLengthFieldOffset(value int) {
	b.BaseReader.SetReadLengthFieldOffset(value)
	b.BaseWriter.SetWriteLengthFieldOffset(value)
}

//SetLengthFieldLength will set the LengthFieldLength for reader/writer
func (b *BaseReadWriteCloser) SetLengthFieldLength(value int) {
	b.BaseReader.SetReadLengthFieldLength(value)
	b.BaseWriter.SetWriteLengthFieldLength(value)
}

//SetLengthAdjustment will set the LengthAdjustment for reader/writer
func (b *BaseReadWriteCloser) SetLengthAdjustment(value int) {
	b.BaseReader.SetReadLengthAdjustment(value)
	b.BaseWriter.SetWriteLengthAdjustment(value)
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
	ByteOrder         binary.ByteOrder
	LengthFieldMagic  int
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment  int
	Buffer            []byte
	Raw               io.Reader
	Timeout           time.Duration
	offset            uint32
	length            uint32
	locker            sync.RWMutex
}

//NewBaseReader will create new Reader by raw reader and buffer size
func NewBaseReader(raw io.Reader, bufferSize int) (reader *BaseReader) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	reader = &BaseReader{
		ByteOrder:         binary.BigEndian,
		LengthFieldMagic:  1,
		LengthFieldLength: 4,
		Buffer:            make([]byte, bufferSize),
		Raw:               raw,
		locker:            sync.RWMutex{},
	}
	return
}

func (b *BaseReader) GetReadByteOrder() (order binary.ByteOrder) {
	order = b.ByteOrder
	return
}

func (b *BaseReader) GetReadLengthFieldMagic() (value int) {
	value = b.LengthFieldMagic
	return
}

func (b *BaseReader) GetReadLengthFieldOffset() (value int) {
	value = b.LengthFieldOffset
	return
}

func (b *BaseReader) GetReadLengthFieldLength() (value int) {
	value = b.LengthFieldLength
	return
}

func (b *BaseReader) GetReadLengthAdjustment() (value int) {
	value = b.LengthAdjustment
	return
}

func (b *BaseReader) SetReadByteOrder(order binary.ByteOrder) {
	b.ByteOrder = order
}

func (b *BaseReader) SetReadLengthFieldMagic(value int) {
	b.LengthFieldMagic = value
}

func (b *BaseReader) SetReadLengthFieldOffset(value int) {
	b.LengthFieldOffset = value
}

func (b *BaseReader) SetReadLengthFieldLength(value int) {
	b.LengthFieldLength = value
}

func (b *BaseReader) SetReadLengthAdjustment(value int) {
	b.LengthAdjustment = value
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

func (b *BaseReader) readFrameLength() (length uint32) {
	for i := 0; i < b.LengthFieldMagic; i++ {
		b.Buffer[b.offset+uint32(b.LengthFieldOffset)+uint32(i)] = 0
	}
	switch b.LengthFieldLength {
	case 1:
		length = uint32(b.Buffer[b.offset+uint32(b.LengthFieldOffset)]) - uint32(b.LengthAdjustment)
	case 2:
		length = uint32(b.ByteOrder.Uint16(b.Buffer[b.offset+uint32(b.LengthFieldOffset):])) - uint32(b.LengthAdjustment)
	case 4:
		length = uint32(b.ByteOrder.Uint32(b.Buffer[b.offset+uint32(b.LengthFieldOffset):])) - uint32(b.LengthAdjustment)
	default:
		panic("not supported LengthFieldLength")
	}
	return
}

//ReadFrame will read raw reader as frame mode. it will return length(4bytes)+data.
//the return []byte is the buffer slice, must be copy to new []byte, it will be change after next read
func (b *BaseReader) ReadFrame() (cmd []byte, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	more := b.length < uint32(b.LengthFieldLength)+1
	for {
		if more {
			err = b.readMore()
			if err != nil {
				break
			}
			if b.length < uint32(b.LengthFieldLength)+1 {
				continue
			}
		}
		frameLength := b.readFrameLength()
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
		more = b.length <= uint32(b.LengthFieldLength)
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
	return xio.RemoteAddr(b.Raw)
}

//BaseWriter implment the frame Writer
type BaseWriter struct {
	ByteOrder         binary.ByteOrder
	LengthFieldMagic  int
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment  int
	Raw               io.Writer
	Timeout           time.Duration
	locker            sync.RWMutex
}

//NewBaseWriter will return new BaseWriter
func NewBaseWriter(raw io.Writer) (writer *BaseWriter) {
	writer = &BaseWriter{
		ByteOrder:         binary.BigEndian,
		LengthFieldMagic:  1,
		LengthFieldLength: 4,
		Raw:               raw,
	}
	return
}

func (b *BaseWriter) GetWriteByteOrder() (order binary.ByteOrder) {
	order = b.ByteOrder
	return
}

func (b *BaseWriter) GetWriteLengthFieldMagic() (value int) {
	value = b.LengthFieldMagic
	return
}

func (b *BaseWriter) GetWriteLengthFieldOffset() (value int) {
	value = b.LengthFieldOffset
	return
}

func (b *BaseWriter) GetWriteLengthFieldLength() (value int) {
	value = b.LengthFieldLength
	return
}

func (b *BaseWriter) GetWriteLengthAdjustment() (value int) {
	value = b.LengthAdjustment
	return
}

func (b *BaseWriter) SetWriteByteOrder(order binary.ByteOrder) {
	b.ByteOrder = order
}

func (b *BaseWriter) SetWriteLengthFieldMagic(value int) {
	b.LengthFieldMagic = value
}

func (b *BaseWriter) SetWriteLengthFieldOffset(value int) {
	b.LengthFieldOffset = value
}

func (b *BaseWriter) SetWriteLengthFieldLength(value int) {
	b.LengthFieldLength = value
}

func (b *BaseWriter) SetWriteLengthAdjustment(value int) {
	b.LengthAdjustment = value
}

//WriteFrame will write data by frame mode, it must have 4 bytes at the begin of buffer to store the frame length.
//genral buffer is (4 bytes)+(user data), 4 bytes will be set the in WriteCmd
func (b *BaseWriter) WriteFrame(buffer []byte) (w int, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	switch b.LengthFieldLength {
	case 1:
		buffer[b.LengthFieldOffset] = byte(len(buffer) + b.LengthAdjustment)
	case 2:
		b.ByteOrder.PutUint16(buffer[b.LengthFieldOffset:], uint16(len(buffer)+b.LengthAdjustment))
	case 4:
		b.ByteOrder.PutUint32(buffer[b.LengthFieldOffset:], uint32(len(buffer)+b.LengthAdjustment))
	default:
		panic("not supported LengthFieldLength")
	}
	for i := 0; i < b.LengthFieldMagic; i++ {
		buffer[b.LengthFieldOffset+i] = byte(rand.Intn(255))
	}
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
	return xio.RemoteAddr(b.Raw)
}

type BasePiper struct {
	Raw               xio.Piper
	BufferSize        int
	Timeout           time.Duration
	ByteOrder         binary.ByteOrder
	LengthFieldMagic  int
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment  int
}

func NewBasePiper(raw xio.Piper, bufferSize int) (piper *BasePiper) {
	piper = &BasePiper{
		Raw:               raw,
		BufferSize:        bufferSize,
		LengthFieldMagic:  1,
		LengthFieldLength: 4,
	}
	return
}

func (b *BasePiper) PipeConn(conn io.ReadWriteCloser, target string) (err error) {
	rwc := NewReadWriteCloser(conn, b.BufferSize)
	rwc.SetTimeout(b.Timeout)
	rwc.SetByteOrder(b.ByteOrder)
	rwc.SetLengthFieldMagic(b.LengthFieldMagic)
	rwc.SetLengthFieldOffset(b.LengthFieldOffset)
	rwc.SetLengthFieldLength(b.LengthFieldLength)
	rwc.SetLengthAdjustment(b.LengthAdjustment)
	err = b.Raw.PipeConn(rwc, target)
	return
}

func (b *BasePiper) Close() (err error) {
	err = b.Raw.Close()
	return
}
