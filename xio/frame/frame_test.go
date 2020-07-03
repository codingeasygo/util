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
)

func TestReadWrite(t *testing.T) {
	//
	{ //one f
		data1 := []byte("one")
		buf := make([]byte, 4+len(data1))
		binary.BigEndian.PutUint32(buf, uint32(4+len(data1)))
		copy(buf[4:], data1)
		raw := bytes.NewBuffer(buf)
		proc := NewReadWriter(raw, 256*1024)
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
			proc := NewReadWriteCloser(r, 256*1024)
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
		proc := NewReadWriter(raw, 256*1024)
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
			proc := NewReadWriter(r, 256*1024)
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
	{ //test frame read write
		r, w, _ := os.Pipe()
		reader := NewReader(r, 1024)
		writer := NewWriter(w)
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
	{ //test close
		NewReadWriteCloser(&net.TCPConn{}, 1024).Close()
	}
	{ //test string
		fmt.Printf("%v\n", NewReader(bytes.NewBuffer(nil), 1024))
		fmt.Printf("%v\n", NewWriter(bytes.NewBuffer(nil)))
	}
}
