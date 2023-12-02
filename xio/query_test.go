package xio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"
)

func TestQueryConn(t *testing.T) {
	query := NewQueryConn()
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	go func() {
		defer waiter.Done()
		io.Copy(query, query)
	}()
	response, err := query.Query(context.Background(), []byte("abc"))
	if err != nil || !bytes.Equal(response, []byte("abc")) {
		t.Error(err)
		return
	}
	query.Close()
	waiter.Wait()

	_, err = query.Query(context.Background(), []byte("abc"))
	if err == nil {
		t.Error(err)
		return
	}

	//not response
	query = NewQueryConn()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = query.Query(ctx, []byte("abc"))
	cancel()
	if err == nil {
		t.Error(err)
		return
	}

	//not request
	query = NewQueryConn()
	query.sendData <- []byte("data")
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = query.Query(ctx, []byte("abc"))
	cancel()
	if err == nil {
		t.Error(err)
		return
	}

	//write twice
	query.recvWait <- 1
	query.Write([]byte("abc"))
	query.Write([]byte("abc"))

	//close twice
	query.Close()
	query.Close()
	query.clearSend()
	query.clearSend()

	//cover
	query.LocalAddr()
	query.RemoteAddr()
	query.SetDeadline(time.Now())
	query.SetReadDeadline(time.Now())
	query.SetWriteDeadline(time.Now())
	query.Network()
	fmt.Println(query.String())
}
