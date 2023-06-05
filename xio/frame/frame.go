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
	//DefaultLengthFieldLength is default frame header length
	DefaultLengthFieldLength = 4
	DefaultBufferSize        = 8 * 1024
)

// ErrFrameTooLarge is the error when the frame head lenght > buffer length
var ErrFrameTooLarge = fmt.Errorf("%v", "frame is too large")

type readDeadlinable interface {
	SetReadDeadline(t time.Time) error
}

type writeDeadlinable interface {
	SetWriteDeadline(t time.Time) error
}

type Header interface {
	GetByteOrder() (order binary.ByteOrder)
	GetLengthFieldMagic() (value int)
	GetLengthFieldOffset() (value int)
	GetLengthFieldLength() (value int)
	GetLengthAdjustment() (value int)
	GetDataOffset() (value int)
	GetDataPrefix() (prefix []byte)
	SetByteOrder(order binary.ByteOrder)
	SetLengthFieldMagic(value int)
	SetLengthFieldOffset(value int)
	SetLengthFieldLength(value int)
	SetLengthAdjustment(value int)
	SetDataOffset(value int)
	SetDataPrefix(prefix []byte)
	WriteHead(buffer []byte)
	ReadHead(buffer []byte) (length uint32)
}

// Reader is interface for read the raw io as frame mode
type Reader interface {
	io.Reader
	Header
	BufferSize() int
	ReadFrame() (frame []byte, err error)
	SetReadTimeout(timeout time.Duration)
	WriteTo(writer io.Writer) (w int64, err error)
}

// Writer is interface for write the raw io as frame mode
type Writer interface {
	io.Writer
	Header
	WriteFrame(buffer []byte) (n int, err error)
	SetWriteTimeout(timeout time.Duration)
	ReadFrom(reader io.Reader) (w int64, err error)
}

// ReadWriter is interface for read/write the raw io as frame mode
type ReadWriter interface {
	Reader
	Writer
	SetTimeout(timeout time.Duration)
}

// ReadWriteCloser is interface for read/write the raw io as frame mode
type ReadWriteCloser interface {
	ReadWriter
	io.Closer
}

type BaseHeader struct {
	ByteOrder         binary.ByteOrder
	LengthFieldMagic  int
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment  int
	DataOffset        int
	DataPrefix        []byte
}

func NewDefaultHeader() (header *BaseHeader) {
	header = &BaseHeader{
		ByteOrder:         binary.BigEndian,
		LengthFieldMagic:  0,
		LengthFieldOffset: 0,
		LengthFieldLength: 4,
		LengthAdjustment:  0,
		DataOffset:        4,
	}
	return
}

func CloneHeader(src Header) (header *BaseHeader) {
	header = &BaseHeader{
		ByteOrder:         src.GetByteOrder(),
		LengthFieldMagic:  src.GetLengthFieldMagic(),
		LengthFieldOffset: src.GetLengthFieldOffset(),
		LengthFieldLength: src.GetLengthFieldLength(),
		LengthAdjustment:  src.GetLengthAdjustment(),
		DataOffset:        src.GetDataOffset(),
		DataPrefix:        src.GetDataPrefix(),
	}
	return
}

func (b *BaseHeader) WriteHead(buffer []byte) {
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
}

func (b *BaseHeader) ReadHead(buffer []byte) (length uint32) {
	for i := 0; i < b.LengthFieldMagic; i++ {
		buffer[uint32(b.LengthFieldOffset)+uint32(i)] = 0
	}
	switch b.LengthFieldLength {
	case 1:
		length = uint32(buffer[uint32(b.LengthFieldOffset)]) - uint32(b.LengthAdjustment)
	case 2:
		length = uint32(b.ByteOrder.Uint16(buffer[+uint32(b.LengthFieldOffset):])) - uint32(b.LengthAdjustment)
	case 4:
		length = uint32(b.ByteOrder.Uint32(buffer[uint32(b.LengthFieldOffset):])) - uint32(b.LengthAdjustment)
	default:
		panic("not supported LengthFieldLength")
	}
	return
}

func (b *BaseHeader) GetByteOrder() (order binary.ByteOrder) {
	order = b.ByteOrder
	return
}

func (b *BaseHeader) GetLengthFieldMagic() (value int) {
	value = b.LengthFieldMagic
	return
}

func (b *BaseHeader) GetLengthFieldOffset() (value int) {
	value = b.LengthFieldOffset
	return
}

func (b *BaseHeader) GetLengthFieldLength() (value int) {
	value = b.LengthFieldLength
	return
}

func (b *BaseHeader) GetLengthAdjustment() (value int) {
	value = b.LengthAdjustment
	return
}

func (b *BaseHeader) GetDataOffset() (value int) {
	value = b.DataOffset
	return
}

func (b *BaseHeader) GetDataPrefix() (prefix []byte) {
	prefix = b.DataPrefix
	return
}

func (b *BaseHeader) SetByteOrder(order binary.ByteOrder) {
	b.ByteOrder = order
}

func (b *BaseHeader) SetLengthFieldMagic(value int) {
	b.LengthFieldMagic = value
}

func (b *BaseHeader) SetLengthFieldOffset(value int) {
	b.LengthFieldOffset = value
}

func (b *BaseHeader) SetLengthFieldLength(value int) {
	b.LengthFieldLength = value
}

func (b *BaseHeader) SetLengthAdjustment(value int) {
	b.LengthAdjustment = value
}

func (b *BaseHeader) SetDataOffset(value int) {
	b.DataOffset = value
}

func (b *BaseHeader) SetDataPrefix(prefix []byte) {
	b.DataPrefix = prefix
}

// NewReader will create new Reader by raw reader and buffer size
func NewReader(raw io.Reader, bufferSize int) (reader *BaseReader) {
	reader = NewBaseReader(raw, bufferSize)
	return
}

// NewWriter will return new BaseWriter
func NewWriter(raw io.Writer) (writer *BaseWriter) {
	writer = NewBaseWriter(raw)
	return
}

// BaseReadWriteCloser is frame reader/writer combiner
type BaseReadWriteCloser struct {
	io.Closer
	Header
	*BaseReader
	*BaseWriter
}

// Close will call the closer
func (b *BaseReadWriteCloser) Close() (err error) {
	if b.Closer != nil {
		err = b.Closer.Close()
	}
	return
}

func (b *BaseReadWriteCloser) String() string {
	return fmt.Sprintf("Reader:%v,Writer:%v", b.BaseReader, b.BaseWriter)
}

// SetTimeout will record the timout
func (b *BaseReadWriteCloser) SetTimeout(timeout time.Duration) {
	b.BaseReader.SetReadTimeout(timeout)
	b.BaseWriter.SetWriteTimeout(timeout)
}

// NewReadWriter will return new ReadWriteCloser
func NewReadWriter(header Header, raw io.ReadWriter, bufferSize int) (frame *BaseReadWriteCloser) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	if header == nil {
		header = NewDefaultHeader()
	} else {
		header = CloneHeader(header)
	}
	closer, _ := raw.(io.Closer)
	frame = &BaseReadWriteCloser{
		Closer:     closer,
		BaseReader: NewBaseReader(raw, bufferSize),
		BaseWriter: NewBaseWriter(raw),
	}
	frame.Header = header
	frame.BaseReader.Header = header
	frame.BaseWriter.Header = header
	return
}

// NewReadWriteCloser will return new ReadWriteCloser
func NewReadWriteCloser(header Header, raw io.ReadWriteCloser, bufferSize int) (frame *BaseReadWriteCloser) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	if header == nil {
		header = NewDefaultHeader()
	} else {
		header = CloneHeader(header)
	}
	frame = &BaseReadWriteCloser{
		Closer:     raw,
		BaseReader: NewBaseReader(raw, bufferSize),
		BaseWriter: NewBaseWriter(raw),
	}
	frame.Header = header
	frame.BaseReader.Header = header
	frame.BaseWriter.Header = header
	return
}

// BaseReader imple read raw connection by frame mode
type BaseReader struct {
	Header
	Buffer  []byte
	Raw     io.Reader
	Timeout time.Duration
	offset  uint32
	length  uint32
	locker  sync.RWMutex
}

// NewBaseReader will create new Reader by raw reader and buffer size
func NewBaseReader(raw io.Reader, bufferSize int) (reader *BaseReader) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	reader = &BaseReader{
		Header: NewDefaultHeader(),
		Buffer: make([]byte, bufferSize),
		Raw:    raw,
		locker: sync.RWMutex{},
	}
	return
}

func (b *BaseReader) BufferSize() int { return len(b.Buffer) }

// readMore will read more data to buffer
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

// ReadFrame will read raw reader as frame mode. it will return length(4bytes)+data.
// the return []byte is the buffer slice, must be copy to new []byte, it will be change after next read
func (b *BaseReader) ReadFrame() (cmd []byte, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	more := b.length < uint32(b.GetLengthFieldLength())+1
	for {
		if more {
			err = b.readMore()
			if err != nil {
				break
			}
			if b.length < uint32(b.GetLengthFieldLength())+1 {
				continue
			}
		}
		frameLength := b.ReadHead(b.Buffer[b.offset : b.offset+b.length])
		if frameLength < 1 {
			err = fmt.Errorf("frame length is zero")
			break
		}
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
		b.offset += frameLength
		b.length -= frameLength
		more = b.length <= uint32(b.GetLengthFieldLength())
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

// Read implment the io.Reader
// it will read the one frame and copy the data to p
func (b *BaseReader) Read(p []byte) (n int, err error) {
	n, err = Read(b, p)
	return
}

// SetReadTimeout will record the timout
func (b *BaseReader) SetReadTimeout(timeout time.Duration) {
	b.Timeout = timeout
}

func (b *BaseReader) WriteTo(writer io.Writer) (w int64, err error) {
	w, err = WriteTo(b, writer)
	return
}

func (b *BaseReader) String() string {
	return xio.RemoteAddr(b.Raw)
}

func Read(reader Reader, p []byte) (n int, err error) {
	offset := reader.GetDataOffset()
	data, err := reader.ReadFrame()
	if err == nil {
		n = copy(p, data[offset:])
	}
	return
}

func WriteTo(reader Reader, writer io.Writer) (w int64, err error) {
	var n int
	var buffer []byte
	offset := reader.GetDataOffset()
	for {
		buffer, err = reader.ReadFrame()
		if err == nil {
			n, err = writer.Write(buffer[offset:])
		}
		if err != nil {
			break
		}
		w += int64(n)
	}
	return
}

// BaseWriter implment the frame Writer
type BaseWriter struct {
	Header
	Raw     io.Writer
	Timeout time.Duration
	locker  sync.RWMutex
}

// NewBaseWriter will return new BaseWriter
func NewBaseWriter(raw io.Writer) (writer *BaseWriter) {
	writer = &BaseWriter{
		Header: NewDefaultHeader(),
		Raw:    raw,
		locker: sync.RWMutex{},
	}
	return
}

// WriteFrame will write data by frame mode, it must have 4 bytes at the begin of buffer to store the frame length.
// genral buffer is (4 bytes)+(user data), 4 bytes will be set the in WriteCmd
func (b *BaseWriter) WriteFrame(buffer []byte) (w int, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	if w, ok := b.Raw.(writeDeadlinable); b.Timeout > 0 && ok {
		w.SetWriteDeadline(time.Now().Add(b.Timeout))
	}
	b.WriteHead(buffer)
	w, err = b.Raw.Write(buffer)
	return
}

// Write implment the io.Writer, the p is user data buffer.
// it will make a new []byte with len(p)+4, the copy data to buffer
func (b *BaseWriter) Write(p []byte) (n int, err error) {
	n, err = Write(b, p)
	return
}

// SetWriteTimeout will record the timout
func (b *BaseWriter) SetWriteTimeout(timeout time.Duration) {
	b.Timeout = timeout
}

func (b *BaseWriter) ReadFrom(reader io.Reader) (w int64, err error) {
	w, err = ReadFrom(b, reader, DefaultBufferSize)
	return
}

func (b *BaseWriter) String() string {
	return xio.RemoteAddr(b.Raw)
}

func Write(writer Writer, p []byte) (n int, err error) {
	offset := writer.GetDataOffset()
	buf := make([]byte, len(p)+offset)
	copy(buf[offset:], p)
	n = len(buf)
	_, err = writer.WriteFrame(buf)
	return
}

func ReadFrom(writer Writer, reader io.Reader, bufferSize int) (w int64, err error) {
	var n int
	buffer := make([]byte, bufferSize)
	offset := writer.GetDataOffset()
	for {
		n, err = reader.Read(buffer[offset:])
		if err == nil {
			n, err = writer.WriteFrame(buffer[:offset+n])
		}
		if err != nil {
			break
		}
		w += int64(n)
	}
	return
}

type BasePiper struct {
	Raw        xio.Piper
	Header     Header
	BufferSize int
	Timeout    time.Duration
}

func NewBasePiper(raw xio.Piper, bufferSize int) (piper *BasePiper) {
	piper = &BasePiper{
		Raw:        raw,
		Header:     NewDefaultHeader(),
		BufferSize: bufferSize,
	}
	return
}

func (b *BasePiper) PipeConn(conn io.ReadWriteCloser, target string) (err error) {
	rwc := NewReadWriteCloser(b.Header, conn, b.BufferSize)
	rwc.SetTimeout(b.Timeout)
	err = b.Raw.PipeConn(rwc, target)
	return
}

func (b *BasePiper) Close() (err error) {
	err = b.Raw.Close()
	return
}

// RawWrapReadWriteCloser is frame reader/writer combiner
type RawWrapReadWriteCloser struct {
	io.Closer
	Header
	*RawWrapReader
	*RawWrapWriter
}

// Close will call the closer
func (r *RawWrapReadWriteCloser) Close() (err error) {
	if r.Closer != nil {
		err = r.Closer.Close()
	}
	return
}

func (r *RawWrapReadWriteCloser) String() string {
	return fmt.Sprintf("Reader:%v,Writer:%v", r.RawWrapReader, r.RawWrapWriter)
}

// SetTimeout will record the timout
func (r *RawWrapReadWriteCloser) SetTimeout(timeout time.Duration) {
	r.RawWrapReader.SetReadTimeout(timeout)
	r.RawWrapWriter.SetWriteTimeout(timeout)
}

// NewRawReadWriter will return new ReadWriteCloser
func NewRawReadWriter(header Header, raw io.ReadWriter, bufferSize int) (frame *RawWrapReadWriteCloser) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	if header == nil {
		header = NewDefaultHeader()
	} else {
		header = CloneHeader(header)
	}
	closer, _ := raw.(io.Closer)
	frame = &RawWrapReadWriteCloser{
		Closer:        closer,
		RawWrapReader: NewRawWrapReader(raw, bufferSize),
		RawWrapWriter: NewRawWrapWriter(raw),
	}
	frame.Header = header
	frame.RawWrapReader.Header = header
	frame.RawWrapWriter.Header = header
	return
}

// NewRawReadWriteCloser will return new ReadWriteCloser
func NewRawReadWriteCloser(header Header, raw io.ReadWriteCloser, bufferSize int) (frame *RawWrapReadWriteCloser) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	if header == nil {
		header = NewDefaultHeader()
	} else {
		header = CloneHeader(header)
	}
	frame = &RawWrapReadWriteCloser{
		Closer:        raw,
		RawWrapReader: NewRawWrapReader(raw, bufferSize),
		RawWrapWriter: NewRawWrapWriter(raw),
	}
	frame.Header = header
	frame.RawWrapReader.Header = header
	frame.RawWrapWriter.Header = header
	return
}

// RawWrapReader imple read raw connection by frame mode
type RawWrapReader struct {
	Header
	Buffer  []byte
	Raw     io.Reader
	Timeout time.Duration
	locker  sync.RWMutex
}

// NewRawWrapReader will create new Reader by raw reader and buffer size
func NewRawWrapReader(raw io.Reader, bufferSize int) (reader *RawWrapReader) {
	if bufferSize < 1 {
		panic("buffer size is < 1")
	}
	reader = &RawWrapReader{
		Header: NewDefaultHeader(),
		Buffer: make([]byte, bufferSize),
		Raw:    raw,
		locker: sync.RWMutex{},
	}
	return
}

func (r *RawWrapReader) BufferSize() int { return len(r.Buffer) }

// ReadFrame will read raw reader as raw mode. it will return DataOffset+data.
func (r *RawWrapReader) ReadFrame() (cmd []byte, err error) {
	r.locker.Lock()
	defer r.locker.Unlock()
	offset := r.GetDataOffset()
	prefix := r.GetDataPrefix()
	l := len(prefix)
	n, err := r.Read(r.Buffer[offset+l:])
	if err != nil {
		return
	}
	cmd = r.Buffer[:offset+l+n]
	r.WriteHead(cmd)
	if l > 0 {
		copy(cmd[offset:offset+l], prefix)
	}
	return
}

// Read implment the io.Reader
// it will read the one frame and copy the data to p
func (r *RawWrapReader) Read(p []byte) (n int, err error) {
	if raw, ok := r.Raw.(readDeadlinable); r.Timeout > 0 && ok {
		raw.SetReadDeadline(time.Now().Add(r.Timeout))
	}
	n, err = r.Raw.Read(p)
	return
}

// SetReadTimeout will record the timout
func (r *RawWrapReader) SetReadTimeout(timeout time.Duration) {
	r.Timeout = timeout
}

func (r *RawWrapReader) WriteTo(writer io.Writer) (w int64, err error) {
	w, err = io.CopyBuffer(writer, r.Raw, r.Buffer)
	return
}

func (r *RawWrapReader) String() string {
	return xio.RemoteAddr(r.Raw)
}

// RawWrapWriter implment the frame Writer
type RawWrapWriter struct {
	Header
	Raw     io.Writer
	Timeout time.Duration
	locker  sync.RWMutex
}

// NewRawWrapWriter will return new RawWrapWriter
func NewRawWrapWriter(raw io.Writer) (writer *RawWrapWriter) {
	writer = &RawWrapWriter{
		Header: NewDefaultHeader(),
		Raw:    raw,
		locker: sync.RWMutex{},
	}
	return
}

func (r *RawWrapWriter) WriteFrame(buffer []byte) (w int, err error) {
	r.locker.Lock()
	defer r.locker.Unlock()
	offset := r.GetDataOffset()
	prefix := r.GetDataPrefix()
	n := offset + len(prefix)
	w, err = r.Write(buffer[n:])
	w += n
	return
}

// Write implment the io.Writer, the p is user data buffer.
// it will make a new []byte with len(p)+4, the copy data to buffer
func (r *RawWrapWriter) Write(p []byte) (n int, err error) {
	if raw, ok := r.Raw.(writeDeadlinable); r.Timeout > 0 && ok {
		raw.SetWriteDeadline(time.Now().Add(r.Timeout))
	}
	n, err = r.Raw.Write(p)
	return
}

// SetWriteTimeout will record the timout
func (r *RawWrapWriter) SetWriteTimeout(timeout time.Duration) {
	r.Timeout = timeout
}

func (r *RawWrapWriter) ReadFrom(reader io.Reader) (w int64, err error) {
	w, err = io.CopyBuffer(r.Raw, reader, make([]byte, DefaultBufferSize))
	return
}

func (r *RawWrapWriter) String() string {
	return xio.RemoteAddr(r.Raw)
}

type WrapWriteCloser struct {
	Header
	io.Writer
	io.Closer
	buffer []byte
	length uint32
}

func NewWrapWriteCloser(next io.WriteCloser, bufferSize int) (writer *WrapWriteCloser) {
	writer = &WrapWriteCloser{
		Header: NewDefaultHeader(),
		Writer: next,
		Closer: next,
		buffer: make([]byte, bufferSize),
	}
	return
}

func NewWrapWriter(next io.Writer, bufferSize int) (writer *WrapWriteCloser) {
	writer = &WrapWriteCloser{
		Header: NewDefaultHeader(),
		Writer: next,
		buffer: make([]byte, bufferSize),
	}
	if closer, ok := next.(io.Closer); ok {
		writer.Closer = closer
	}
	return
}

func (w *WrapWriteCloser) readFrame(header uint32, buffer []byte) (size uint32, err error) {
	bufSize := uint32(len(buffer))
	if bufSize < header {
		return
	}
	frameLength := w.ReadHead(buffer)
	if frameLength > uint32(len(w.buffer)) {
		err = ErrFrameTooLarge
		return
	}
	if bufSize >= frameLength {
		size = frameLength
	}
	return
}

func (w *WrapWriteCloser) Write(p []byte) (writed int, err error) {
	recvSize := uint32(len(p))
	recvBuf := p
	header := uint32(w.GetLengthFieldLength())
	offset := uint32(w.GetDataOffset())
	frameSize := uint32(0)
	n := 0
	for {
		if w.length < 1 {
			frameSize, err = w.readFrame(header, recvBuf)
			if err != nil {
				break
			}
			if frameSize < 1 { //need more data
				n = copy(w.buffer, recvBuf)
				writed += n
				w.length += uint32(n)
				break
			}
			n, err = w.Writer.Write(recvBuf[offset:frameSize])
			if err != nil {
				break
			}
			writed += n
			recvBuf = recvBuf[frameSize:]
			recvSize -= frameSize
			if recvSize < 1 {
				break
			}
		} else {
			if len(recvBuf) > 0 {
				n = copy(w.buffer[w.length:], recvBuf)
				writed += n
				recvBuf = recvBuf[n:]
				recvSize -= uint32(n)
				w.length += uint32(n)
			}
			frameSize, err = w.readFrame(header, w.buffer[:w.length])
			if err != nil {
				break
			}
			if frameSize < 1 { //need more data
				break
			}
			n, err = w.Writer.Write(w.buffer[offset:frameSize])
			if err != nil {
				break
			}
			writed += n
			copy(w.buffer[0:], w.buffer[frameSize:w.length])
			w.length -= frameSize
		}
	}
	return
}

func (w *WrapWriteCloser) Close() (err error) {
	if w.Closer != nil {
		err = w.Closer.Close()
	}
	return
}

type WrapReadCloser struct {
	Header
	io.Reader
	io.Closer
}

func NewWrapReadCloser(from io.ReadCloser) (reader *WrapReadCloser) {
	reader = &WrapReadCloser{
		Header: NewDefaultHeader(),
		Reader: from,
		Closer: from,
	}
	return
}

func NewWrapReader(from io.Reader) (reader *WrapReadCloser) {
	reader = &WrapReadCloser{
		Header: NewDefaultHeader(),
		Reader: from,
	}
	if closer, ok := from.(io.Closer); ok {
		reader.Closer = closer
	}
	return
}

func (w *WrapReadCloser) Read(p []byte) (n int, err error) {
	offset := w.GetDataOffset()
	n, err = w.Reader.Read(p[offset:])
	if err == nil {
		n += offset
		w.WriteHead(p[:n])
	}
	return
}

func (w *WrapReadCloser) Close() (err error) {
	if w.Closer != nil {
		err = w.Closer.Close()
	}
	return
}
