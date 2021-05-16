package xio

import (
	"net"
	"sync"
	"time"
)

type TimeoutListener struct {
	net.Listener
	Delay   time.Duration
	Timeout time.Duration
	waiter  sync.WaitGroup
	running bool
	allConn map[*timeoutConn]int
	allLock sync.RWMutex
}

func NewTimeoutListener(ln net.Listener, timeout time.Duration) (listener *TimeoutListener) {
	listener = &TimeoutListener{
		Listener: ln,
		Delay:    time.Second,
		Timeout:  timeout,
		waiter:   sync.WaitGroup{},
		running:  true,
		allConn:  map[*timeoutConn]int{},
		allLock:  sync.RWMutex{},
	}
	listener.waiter.Add(1)
	go listener.runTimeout()
	return
}

func (t *TimeoutListener) runTimeout() {
	for t.running {
		now := time.Now()
		t.allLock.Lock()
		for c := range t.allConn {
			if now.Sub(c.latest) >= t.Timeout {
				c.rawClose()
				delete(t.allConn, c)
			}
		}
		t.allLock.Unlock()
		time.Sleep(t.Delay)
	}
	t.waiter.Done()
}

func (t *TimeoutListener) Accept() (conn net.Conn, err error) {
	conn, err = t.Listener.Accept()
	if err == nil {
		t.allLock.Lock()
		c := &timeoutConn{Conn: conn, listener: t, latest: time.Now()}
		t.allConn[c] = 1
		t.allLock.Unlock()
		conn = c
	}
	return
}

func (t *TimeoutListener) Close() (err error) {
	err = t.Listener.Close()
	t.running = false
	t.waiter.Wait()
	return
}

func (t *TimeoutListener) closeConn(c *timeoutConn) {
	t.allLock.Lock()
	delete(t.allConn, c)
	t.allLock.Unlock()
}

type timeoutConn struct {
	net.Conn
	listener *TimeoutListener
	latest   time.Time
}

func (t *timeoutConn) Read(p []byte) (n int, err error) {
	t.latest = time.Now()
	n, err = t.Conn.Read(p)
	return
}

func (t *timeoutConn) Write(p []byte) (n int, err error) {
	t.latest = time.Now()
	n, err = t.Conn.Write(p)
	return
}

func (t *timeoutConn) rawClose() (err error) {
	t.latest = time.Now()
	err = t.Conn.Close()
	return
}

func (t *timeoutConn) Close() (err error) {
	err = t.rawClose()
	t.listener.closeConn(t)
	return
}
