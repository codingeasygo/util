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
	//
	{ //one frame 8
		data1 := []byte("one")
		buf := make([]byte, 1+len(data1))
		buf[0] = byte(len(data1) + 1)
		copy(buf[1:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(&deadlineRWC{ReadWriter: raw}, 256*1024)
		proc.SetTimeout(time.Second)
		proc.SetLengthFieldLength(1)
		proc.SetLengthFieldMagic(0)
		proc.SetLengthFieldOffset(0)
		proc.SetLengthAdjustment(0)
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
	{ //one frame 16
		data1 := []byte("one")
		buf := make([]byte, 2+len(data1))
		binary.BigEndian.PutUint16(buf, uint16(2+len(data1)))
		copy(buf[2:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(&deadlineRWC{ReadWriter: raw}, 256*1024)
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
	{ //one frame 32
		data1 := []byte("one")
		buf := make([]byte, 4+len(data1))
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(&deadlineRWC{ReadWriter: raw}, 256*1024)
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
	{ //one frame splice
		data1 := []byte("one")
		r, w, _ := os.Pipe()
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			proc := NewReadWriteCloser(&deadlineRWC{ReadWriter: r}, 256*1024)
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
	{ //two frame
		data1 := []byte("two1")
		data2 := []byte("two2")
		buf := make([]byte, 8+len(data1)+len(data2))
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		binary.BigEndian.PutUint32(buf[4+len(data1):], uint32(4+len(data2)))
		copy(buf[8+len(data1):], data2)
		raw := bytes.NewBuffer(buf)
		//
		proc := NewReadWriter(&deadlineRWC{ReadWriter: raw}, 256*1024)
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
	{ //two frame splice
		data1 := []byte("splice1")
		data2 := []byte("splice2")
		r, w, _ := os.Pipe()
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			proc := NewReadWriter(&deadlineRWC{ReadWriter: r}, 256*1024)
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
	{ //two frame splice 2
		data1 := []byte("splice1")
		data2 := []byte("splice2")
		r, w, _ := os.Pipe()
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			proc := NewReadWriter(&deadlineRWC{ReadWriter: r}, 256*1024)
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
	{ //test frame read write 8
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		reader.LengthFieldLength = 1
		reader.LengthFieldMagic = 0
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
		writer.SetWriteTimeout(time.Second)
		writer.LengthFieldLength = 1
		writer.LengthFieldMagic = 0
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
	{ //test frame read write 16
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		reader.LengthFieldLength = 2
		reader.LengthFieldMagic = 0
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
		writer.SetWriteTimeout(time.Second)
		writer.LengthFieldLength = 2
		writer.LengthFieldMagic = 0
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
	{ //test frame read write
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		writer := NewWriter(&deadlineRWC{ReadWriter: w})
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
	{ //test too large
		buf := make([]byte, 1024)
		binary.BigEndian.PutUint32(buf, 1000000)
		proc := NewReadWriter(bytes.NewBuffer(buf), 1024)
		_, err := proc.ReadFrame()
		if err == nil {
			t.Error(err)
			return
		}
	}
	{
		//
		NewReader(nil, 1024).GetReadByteOrder()
		NewReader(nil, 1024).SetReadByteOrder(binary.BigEndian)
		NewReader(nil, 1024).GetReadLengthFieldMagic()
		NewReader(nil, 1024).SetReadLengthFieldMagic(1)
		NewReader(nil, 1024).GetReadLengthFieldOffset()
		NewReader(nil, 1024).SetReadLengthFieldOffset(1)
		NewReader(nil, 1024).GetReadLengthFieldLength()
		NewReader(nil, 1024).SetReadLengthFieldLength(1)
		NewReader(nil, 1024).GetReadLengthAdjustment()
		NewReader(nil, 1024).SetReadLengthAdjustment(1)
		//
		NewWriter(nil).GetWriteByteOrder()
		NewWriter(nil).SetWriteByteOrder(binary.BigEndian)
		NewWriter(nil).GetWriteLengthFieldMagic()
		NewWriter(nil).SetWriteLengthFieldMagic(1)
		NewWriter(nil).GetWriteLengthFieldOffset()
		NewWriter(nil).SetWriteLengthFieldOffset(1)
		NewWriter(nil).GetWriteLengthFieldLength()
		NewWriter(nil).SetWriteLengthFieldLength(1)
		NewWriter(nil).GetWriteLengthAdjustment()
		NewWriter(nil).SetWriteLengthAdjustment(1)
		//
		NewReadWriter(nil, 1024).GetByteOrder()
		NewReadWriter(nil, 1024).SetByteOrder(binary.BigEndian)
		NewReadWriter(nil, 1024).GetLengthFieldMagic()
		NewReadWriter(nil, 1024).SetLengthFieldMagic(1)
		NewReadWriter(nil, 1024).GetLengthFieldOffset()
		NewReadWriter(nil, 1024).SetLengthFieldOffset(1)
		NewReadWriter(nil, 1024).GetLengthFieldLength()
		NewReadWriter(nil, 1024).SetLengthFieldLength(1)
		NewReadWriter(nil, 1024).GetLengthAdjustment()
		NewReadWriter(nil, 1024).SetLengthAdjustment(1)
	}
	{ //test close
		NewReadWriteCloser(&net.TCPConn{}, 1024).Close()
	}
	{ //test string
		fmt.Printf("%v\n", NewReader(bytes.NewBuffer(nil), 1024))
		fmt.Printf("%v\n", NewWriter(bytes.NewBuffer(nil)))
		fmt.Printf("%v\n", NewReadWriteCloser(nil, 1024))
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
		NewReadWriter(nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		NewReadWriteCloser(nil, -1)
	}()
	func() {
		defer func() {
			recover()
		}()
		rwc := NewReadWriteCloser(nil, 1024)
		rwc.SetLengthFieldLength(10)
		rwc.WriteFrame(nil)
	}()
	func() {
		defer func() {
			recover()
		}()
		rwc := NewReadWriteCloser(nil, 1024)
		rwc.SetLengthFieldLength(10)
		rwc.readFrameLength()
	}()
}
