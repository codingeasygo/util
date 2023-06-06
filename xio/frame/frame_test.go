package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/codingeasygo/util/xdebug"
	"github.com/codingeasygo/util/xio"
)

type deadlineRWC struct {
	io.ReadWriter
}

func (d *deadlineRWC) SetReadDeadline(t time.Time) error {
	return nil
}

func (d *deadlineRWC) SetWriteDeadline(t time.Time) error {
	return nil
}

func (d *deadlineRWC) Close() (err error) {
	if c, ok := d.ReadWriter.(io.Closer); ok {
		err = c.Close()
	}
	return
}

func TestReadWrite(t *testing.T) {
	tester := xdebug.CaseTester{
		0:  1,
		11: 1,
	}
	//
	if tester.Run() { //one frame 8
		data1 := []byte("one")
		buf := make([]byte, 1+len(data1))
		buf[0] = byte(len(data1) + 1)
		copy(buf[1:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(nil, &deadlineRWC{ReadWriter: raw}, 256*1024)
		proc.SetTimeout(time.Second)
		proc.SetLengthFieldLength(1)
		proc.SetDataOffset(1)
		f, err := proc.ReadFrame()
		if err != nil || !bytes.Equal(f[1:], data1) {
			t.Error(err)
			return
		}
		_, err = proc.ReadFrame()
		if err != io.EOF {
			t.Error(err)
			return
		}
	}
	if tester.Run() { //one frame 16
		data1 := []byte("one")
		buf := make([]byte, 2+len(data1))
		binary.BigEndian.PutUint16(buf, uint16(2+len(data1)))
		copy(buf[2:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(nil, &deadlineRWC{ReadWriter: raw}, 256*1024)
		proc.SetTimeout(time.Second)
		proc.SetLengthFieldLength(2)
		proc.SetLengthFieldMagic(0)
		f, err := proc.ReadFrame()
		if err != nil || !bytes.Equal(f[2:], data1) {
			t.Error(err)
			return
		}
		_, err = proc.ReadFrame()
		if err != io.EOF {
			t.Error(err)
			return
		}
	}
	if tester.Run() { //one frame 32
		data1 := []byte("one")
		buf := make([]byte, 4+len(data1))
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(nil, &deadlineRWC{ReadWriter: raw}, 256*1024)
		proc.SetLengthFieldMagic(1)
		proc.SetTimeout(time.Second)
		f, err := proc.ReadFrame()
		if err != nil || !bytes.Equal(f[4:], data1) {
			t.Error(err)
			return
		}
		_, err = proc.ReadFrame()
		if err != io.EOF {
			t.Error(err)
			return
		}
	}
	if tester.Run() { //one frame splice
		data1 := []byte("one")
		r, w, _ := os.Pipe()
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			proc := NewReadWriteCloser(nil, &deadlineRWC{ReadWriter: r}, 256*1024)
			proc.SetLengthFieldMagic(1)
			proc.SetTimeout(time.Second)
			f, err := proc.ReadFrame()
			if err != nil || !bytes.Equal(f[4:], data1) {
				t.Error(err)
				return
			}
			_, err = proc.ReadFrame()
			if err != io.EOF {
				t.Error(err)
				return
			}
			wait.Done()
		}()
		buf := make([]byte, uint32(4+len(data1)))
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		w.Write(buf[0:3])
		time.Sleep(time.Millisecond)
		w.Write(buf[3:])
		time.Sleep(time.Millisecond)
		w.Close()
		time.Sleep(time.Millisecond)
		wait.Wait()
	}
	//
	if tester.Run() { //two frame
		data1 := []byte("two1")
		data2 := []byte("two2")
		buf := make([]byte, 8+len(data1)+len(data2))
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		binary.BigEndian.PutUint32(buf[4+len(data1):], uint32(4+len(data2)))
		copy(buf[8+len(data1):], data2)
		raw := bytes.NewBuffer(buf)
		//
		proc := NewReadWriter(nil, &deadlineRWC{ReadWriter: raw}, 256*1024)
		proc.SetLengthFieldMagic(1)
		proc.SetTimeout(time.Second)
		f, err := proc.ReadFrame()
		if err != nil || !bytes.Equal(f[4:], data1) {
			t.Error(err)
			return
		}
		f, err = proc.ReadFrame()
		if err != nil || !bytes.Equal(f[4:], data2) {
			t.Error(err)
			return
		}
		_, err = proc.ReadFrame()
		if err != io.EOF {
			t.Error(err)
			return
		}
	}
	//
	if tester.Run() { //two frame splice
		data1 := []byte("splice1")
		data2 := []byte("splice2")
		r, w, _ := os.Pipe()
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			proc := NewReadWriter(nil, &deadlineRWC{ReadWriter: r}, 256*1024)
			proc.SetLengthFieldMagic(1)
			proc.SetTimeout(time.Second)
			f, err := proc.ReadFrame()
			if err != nil || !bytes.Equal(f[4:], data1) {
				t.Error(err)
				return
			}
			f, err = proc.ReadFrame()
			if err != nil || !bytes.Equal(f[4:], data2) {
				t.Error(err)
				return
			}
			_, err = proc.ReadFrame()
			if err != io.EOF {
				t.Error(err)
				return
			}
			wait.Done()
		}()
		buf := make([]byte, 1024)
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		binary.BigEndian.PutUint32(buf[4+len(data1):], uint32(4+len(data2)))
		copy(buf[8+len(data1):], data2[:1])
		w.Write(buf[:8+len(data1)+1])
		time.Sleep(time.Millisecond)
		w.Write(data2[1:])
		time.Sleep(time.Millisecond)
		w.Close()
		time.Sleep(time.Millisecond)
		wait.Wait()
	}
	//
	if tester.Run() { //two frame splice 2
		data1 := []byte("splice1")
		data2 := []byte("splice2")
		r, w, _ := os.Pipe()
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			proc := NewReadWriter(nil, &deadlineRWC{ReadWriter: r}, 256*1024)
			proc.SetLengthFieldMagic(1)
			proc.SetTimeout(time.Second)
			f, err := proc.ReadFrame()
			if err != nil || !bytes.Equal(f[4:], data1) {
				t.Error(err)
				return
			}
			f, err = proc.ReadFrame()
			if err != nil || !bytes.Equal(f[4:], data2) {
				t.Error(err)
				return
			}
			_, err = proc.ReadFrame()
			if err != io.EOF {
				t.Error(err)
				return
			}
			wait.Done()
		}()
		buf := make([]byte, 1024)
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		binary.BigEndian.PutUint32(buf[4+len(data1):], uint32(4+len(data2)))
		copy(buf[8+len(data1):], data2)
		w.Write(buf[:4+len(data1)+3])
		time.Sleep(time.Millisecond)
		w.Write(buf[4+len(data1)+3 : 4+len(data1)+4+len(data2)])
		time.Sleep(time.Millisecond)
		w.Close()
		time.Sleep(time.Millisecond)
		wait.Wait()
	}
	if tester.Run() { //test frame read write 8
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		reader.SetLengthFieldLength(1)
		reader.SetLengthFieldMagic(0)
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
		writer.SetWriteTimeout(time.Second)
		writer.SetLengthFieldLength(1)
		writer.SetLengthFieldMagic(0)
		readed := bytes.NewBuffer(nil)
		waiter := make(chan int, 1)
		go func() {
			io.Copy(readed, reader)
			waiter <- 1
		}()
		writed := bytes.NewBuffer(nil)
		count := rand.Intn(10) + 1
		for i := 0; i < count; i++ {
			fmt.Fprintf(writer, "data-%v\n", i)
			fmt.Fprintf(writed, "data-%v\n", i)
		}
		w.Close()
		<-waiter
		if !bytes.Equal(readed.Bytes(), writed.Bytes()) {
			fmt.Printf("readed:\n%v\n", (readed.Bytes()))
			fmt.Printf("writed:\n%v\n", (writed.Bytes()))
			t.Error("error")
			return
		}
	}
	if tester.Run() { //test frame read write 16
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		reader.SetLengthFieldLength(2)
		reader.SetLengthFieldMagic(0)
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
		writer.SetWriteTimeout(time.Second)
		writer.SetLengthFieldLength(2)
		writer.SetLengthFieldMagic(0)
		readed := bytes.NewBuffer(nil)
		waiter := make(chan int, 1)
		go func() {
			io.Copy(readed, reader)
			waiter <- 1
		}()
		writed := bytes.NewBuffer(nil)
		count := rand.Intn(10) + 1
		for i := 0; i < count; i++ {
			fmt.Fprintf(writer, "data-%v\n", i)
			fmt.Fprintf(writed, "data-%v\n", i)
		}
		w.Close()
		<-waiter
		if !bytes.Equal(readed.Bytes(), writed.Bytes()) {
			fmt.Printf("readed:\n%v\n", (readed.Bytes()))
			fmt.Printf("writed:\n%v\n", (writed.Bytes()))
			t.Error("error")
			return
		}
	}
	if tester.Run() { //test copy read
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		reader.SetLengthFieldMagic(1)
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
		writer.SetLengthFieldMagic(1)
		writer.SetWriteTimeout(time.Second)
		readed := bytes.NewBuffer(nil)
		waiter := make(chan int, 1)
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := reader.Read(buf)
				if err != nil {
					break
				}
				readed.Write(buf[0:n])
			}
			io.Copy(readed, reader)
			waiter <- 1
		}()
		writed := bytes.NewBuffer(nil)
		count := rand.Intn(10) + 1
		for i := 0; i < count; i++ {
			fmt.Fprintf(writer, "data-%v\n", i)
			fmt.Fprintf(writed, "data-%v\n", i)
		}
		w.Close()
		<-waiter
		if !bytes.Equal(readed.Bytes(), writed.Bytes()) {
			fmt.Printf("readed:\n%v\n", (readed.Bytes()))
			fmt.Printf("writed:\n%v\n", (writed.Bytes()))
			t.Error("error")
			return
		}
	}
	if tester.Run() { //test copy write to
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		reader.SetLengthFieldMagic(1)
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
		writer.SetLengthFieldMagic(1)
		writer.SetWriteTimeout(time.Second)
		readed := bytes.NewBuffer(nil)
		waiter := make(chan int, 1)
		go func() {
			io.Copy(readed, reader)
			waiter <- 1
		}()
		writed := bytes.NewBuffer(nil)
		count := rand.Intn(10) + 1
		for i := 0; i < count; i++ {
			fmt.Fprintf(writer, "data-%v\n", i)
			fmt.Fprintf(writed, "data-%v\n", i)
		}
		w.Close()
		<-waiter
		if !bytes.Equal(readed.Bytes(), writed.Bytes()) {
			fmt.Printf("readed:\n%v\n", (readed.Bytes()))
			fmt.Printf("writed:\n%v\n", (writed.Bytes()))
			t.Error("error")
			return
		}
	}
	if tester.Run() { //test copy write
		cr, cw, _ := os.Pipe()
		sr, sw, _ := os.Pipe()
		reader := NewReader(cr, 1024)
		reader.SetLengthFieldMagic(1)
		writer := NewWriter(&deadlineRWC{ReadWriter: cw})
		writer.SetLengthFieldMagic(1)
		writer.SetWriteTimeout(time.Second)
		readed := bytes.NewBuffer(nil)
		waiter := make(chan int, 1)
		go func() {
			io.Copy(writer, sr)
			waiter <- 1
		}()
		go func() {
			io.Copy(readed, reader)
			waiter <- 1
		}()
		writed := bytes.NewBuffer(nil)
		count := rand.Intn(10) + 1
		for i := 0; i < count; i++ {
			fmt.Fprintf(sw, "data-%v\n", i)
			fmt.Fprintf(writed, "data-%v\n", i)
		}
		time.Sleep(500 * time.Millisecond)
		sw.Close()
		<-waiter
		if !bytes.Equal(readed.Bytes(), writed.Bytes()) {
			fmt.Printf("readed:\n%v\n", (readed.Bytes()))
			fmt.Printf("writed:\n%v\n", (writed.Bytes()))
			t.Error("error")
			return
		}
	}
	if tester.Run() { //test too large
		buf := make([]byte, 1024)
		binary.BigEndian.PutUint32(buf, 1000000)
		proc := NewReadWriter(nil, bytes.NewBuffer(buf), 1024)
		_, err := proc.ReadFrame()
		if err == nil {
			t.Error(err)
			return
		}
	}
	if tester.Run() { //test frame header invalid
		buf := make([]byte, 1024)
		proc := NewReadWriter(nil, bytes.NewBuffer(buf), 1024)
		_, err := proc.ReadFrame()
		if err == nil {
			t.Error(err)
			return
		}
	}
	if tester.Run() { //for cover
		NewReadWriter(nil, nil, 1024).GetByteOrder()
		NewReadWriter(nil, nil, 1024).SetByteOrder(binary.BigEndian)
		NewReadWriter(nil, nil, 1024).GetLengthFieldMagic()
		NewReadWriter(nil, nil, 1024).SetLengthFieldMagic(1)
		NewReadWriter(nil, nil, 1024).GetLengthFieldOffset()
		NewReadWriter(nil, nil, 1024).SetLengthFieldOffset(1)
		NewReadWriter(nil, nil, 1024).GetLengthFieldLength()
		NewReadWriter(nil, nil, 1024).SetLengthFieldLength(1)
		NewReadWriter(nil, nil, 1024).GetLengthAdjustment()
		NewReadWriter(nil, nil, 1024).SetLengthAdjustment(1)
		NewReadWriter(nil, nil, 1024).BufferSize()
	}
	if tester.Run() { //test close
		NewReadWriteCloser(nil, &net.TCPConn{}, 1024).Close()
	}
	if tester.Run() { //test string
		fmt.Printf("%v\n", NewReader(bytes.NewBuffer(nil), 1024))
		fmt.Printf("%v\n", NewWriter(bytes.NewBuffer(nil)))
		fmt.Printf("%v\n", NewReadWriteCloser(nil, nil, 1024))
	}
}

func TestPiper(t *testing.T) {
	piper := NewBasePiper(xio.PiperF(func(conn io.ReadWriteCloser, target string) (err error) {
		err = fmt.Errorf("error")
		return
	}), 1024)
	piper.PipeConn(nil, "target")
	piper.Close()
}

func TestRaw(t *testing.T) {
	tester := xdebug.CaseTester{
		0: 0,
		3: 1,
	}
	if tester.Run() { //read/write frame
		r, w, _ := os.Pipe()
		reader := NewReader(r, DefaultBufferSize)
		writer := NewWriter(w)
		go func() {
			buffer := bytes.NewBufferString("abc")
			src := NewRawReadWriter(nil, &deadlineRWC{ReadWriter: buffer}, DefaultBufferSize)
			src.SetTimeout(time.Second)
			frame, err := src.ReadFrame()
			if err != nil {
				t.Error(err)
				return
			}
			_, err = writer.WriteFrame(frame)
			if err != nil {
				t.Error(err)
				return
			}
			src.Close()
		}()
		buffer := bytes.NewBuffer(nil)
		dst := NewRawReadWriteCloser(nil, &deadlineRWC{ReadWriter: buffer}, DefaultBufferSize)
		dst.SetTimeout(time.Second)
		frame, err := reader.ReadFrame()
		if err != nil {
			t.Error(err)
			return
		}
		_, err = dst.WriteFrame(frame)
		if err != nil {
			t.Error(err)
			return
		}
		if buffer.String() != "abc" {
			t.Error("error")
			return
		}
		dst.Close()
	}
	if tester.Run() { //read/write from
		src := NewRawReadWriteCloser(nil, &deadlineRWC{ReadWriter: bytes.NewBufferString("abc")}, DefaultBufferSize)
		src.SetReadTimeout(time.Second)
		src.WriteTo(os.Stdout)

		dst := NewRawWrapWriter(&deadlineRWC{ReadWriter: os.Stdout})
		dst.SetWriteTimeout(time.Second)
		dst.ReadFrom(bytes.NewBufferString("abc"))

		fmt.Printf("src->%v\n", src)
		fmt.Printf("dst->%v\n", dst)

		src.BufferSize()
	}
	if tester.Run() { //read placeholder
		src := NewRawReadWriteCloser(nil, &deadlineRWC{ReadWriter: bytes.NewBufferString("abc")}, DefaultBufferSize)
		src.SetDataPrefix([]byte("123"))
		frame, err := src.ReadFrame()
		if err != nil || string(frame[src.GetDataOffset():]) != "123abc" {
			t.Errorf("%v", err)
			return
		}
	}
}

func TestPass(t *testing.T) {
	tester := xdebug.CaseTester{
		0: 1,
		9: 1,
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		writer := NewBaseWriter(wrapper)
		n, err := writer.Write([]byte("abc"))
		if err != nil || wrapper.length > 0 || buffer.String() != "abc" {
			t.Errorf("%v,%v", err, n)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0})
		wrapper.Write([]byte{0, 0, 7})
		wrapper.Write([]byte{97, 98, 99})
		if wrapper.length > 0 || buffer.String() != "abc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0})
		wrapper.Write([]byte{0, 0})
		wrapper.Write([]byte{7})
		wrapper.Write([]byte{97, 98, 99})
		if wrapper.length > 0 || buffer.String() != "abc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0, 0, 0, 7})
		wrapper.Write([]byte{97, 98, 99})
		if wrapper.length > 0 || buffer.String() != "abc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0, 0, 0, 7, 97})
		wrapper.Write([]byte{98, 99})
		if wrapper.length > 0 || buffer.String() != "abc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0, 0, 0, 7, 97})
		wrapper.Write([]byte{98})
		wrapper.Write([]byte{99})
		if wrapper.length > 0 || buffer.String() != "abc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0, 0, 0, 7, 97, 98, 99, 0, 0})
		wrapper.Write([]byte{0, 7, 97, 98, 99})
		if wrapper.length > 0 || buffer.String() != "abcabc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{0, 0, 0, 7, 97, 98, 99, 0, 0, 0, 7, 97, 98, 99})
		if wrapper.length > 0 || buffer.String() != "abcabc" {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := make([]byte, 1024)
		data := bytes.NewBuffer([]byte("abc"))
		wrapper := NewPassReader(data)
		reader := NewBaseReader(wrapper, 1024)
		n, err := reader.Read(buffer)
		if err != nil || string(buffer[0:n]) != "abc" {
			t.Errorf("%v,%v", err, n)
			return
		}
		wrapper.Close()
		NewPassReadCloser(NewRawReadWriteCloser(nil, nil, 1024))
		wrapper = NewPassReader(NewRawReadWriteCloser(nil, nil, 1024))
		wrapper.Close()
	}
	//
	//test error
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		_, err := wrapper.Write([]byte{255, 255, 255, 7, 97, 98, 99})
		if err != ErrFrameTooLarge {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
	if tester.Run() {
		buffer := bytes.NewBuffer(nil)
		wrapper := NewPassWriter(buffer, 1024)
		wrapper.Write([]byte{255, 255})
		_, err := wrapper.Write([]byte{255, 7, 97, 98, 99})
		if err != ErrFrameTooLarge {
			t.Errorf("%v,%v", wrapper.length, buffer)
			return
		}
	}
}

type errrWriter struct {
}

func (e *errrWriter) Write(p []byte) (n int, err error) {
	err = fmt.Errorf("error")
	return
}

func (e *errrWriter) Close() (err error) {
	return
}

func TestError(t *testing.T) {
	func() {
		defer func() {
			recover()
		}()
		NewBaseReader(nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		NewReadWriter(nil, nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		NewReadWriteCloser(nil, nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		rwc := NewReadWriteCloser(nil, nil, 1024)
		rwc.SetLengthFieldLength(10)
		rwc.WriteFrame(nil)
	}()
	func() {
		defer func() {
			recover()
		}()
		rwc := NewReadWriteCloser(nil, nil, 1024)
		rwc.SetLengthFieldLength(10)
		rwc.ReadHead([]byte("12222"))
	}()
	func() {
		defer func() {
			recover()
		}()
		rwc := NewReadWriteCloser(nil, nil, 1024)
		rwc.SetLengthFieldLength(10)
		rwc.WriteHead([]byte("12222"))
	}()
	func() {
		defer func() {
			recover()
		}()
		NewRawWrapReader(nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		NewRawReadWriter(nil, nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		NewRawReadWriteCloser(nil, nil, -1)
	}()
	func() {
		NewPassWriter(&errrWriter{}, 1024)
		wrapper := NewPassWriteCloser(&errrWriter{}, 1024)
		wrapper.Write([]byte{0, 0, 0, 7, 97, 98, 99})
		wrapper.Write([]byte{0, 0, 0, 7})
		wrapper.Write([]byte{97, 98, 99})
		wrapper.Close()
	}()
}

func TestEqual(t *testing.T) {
	var a = bytes.NewBuffer(nil)
	var b = bytes.NewBuffer(nil)
	var r io.Reader
	var w io.Writer

	r, w = a, a
	fmt.Printf("--->%v\n", interface{}(r) == interface{}(w))

	r, w = a, b
	fmt.Printf("--->%v\n", interface{}(r) == interface{}(w))
}
