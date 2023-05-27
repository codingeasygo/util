package proxy

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/codingeasygo/util/proxy/socks"
	"github.com/codingeasygo/util/xdebug"
	"github.com/codingeasygo/util/xio"
)

type RouterPiperDialer struct {
	Router string
	Next   xio.PiperDialer
}

func (r *RouterPiperDialer) DialPiper(uri string, bufferSize int) (raw xio.Piper, err error) {
	raw, err = r.Next.DialPiper(strings.Replace(r.Router, "${HOST}", uri, -1), bufferSize)
	return
}

type Forward struct {
	Name       string
	BufferSize int
	Dialer     xio.PiperDialer
	forwardLck sync.RWMutex
	forwardAll map[string][]interface{}
}

func NewForward(name string) (forward *Forward) {
	forward = &Forward{
		Name:       name,
		BufferSize: 8 * 1024,
		Dialer:     xio.PiperDialerF(xio.DialNetPiper),
		forwardLck: sync.RWMutex{},
		forwardAll: map[string][]interface{}{},
	}
	return
}

// StartForward will forward address to uri
func (f *Forward) StartForward(name string, listen *url.URL, router string) (listener net.Listener, err error) {
	f.forwardLck.Lock()
	defer f.forwardLck.Unlock()
	if f.forwardAll[name] != nil || len(name) < 1 {
		err = fmt.Errorf("the name(%v) is already used", name)
		WarnLog("Forward(%v) start forward by %v fail with %v", f.Name, listen, router, err)
		return
	}
	switch listen.Scheme {
	case "socks":
		sp := socks.NewServer()
		sp.BufferSize = f.BufferSize
		sp.Dialer = &RouterPiperDialer{Router: router, Next: f.Dialer}
		listener, err = sp.Start(listen.Host)
		if err == nil {
			f.forwardAll[name] = []interface{}{listen.Scheme, listener, listen}
			InfoLog("Forward(%v) start socket forward on %v success by %v->%v", f.Name, listener.Addr(), listen, router)
		}
	case "proxy":
		dialer := &RouterPiperDialer{Router: router, Next: f.Dialer}
		sp := NewServer(dialer)
		sp.SOCKS.BufferSize = f.BufferSize
		sp.HTTP.BufferSize = f.BufferSize
		listener, err = sp.Start(listen.Host)
		if err == nil {
			f.forwardAll[name] = []interface{}{listen.Scheme, listener, listen, router}
			InfoLog("Forward(%v) start proxy forward on %v success by %v->%v", f.Name, listener.Addr(), listen, router)
		}
	default:
		listener, err = net.Listen(listen.Scheme, listen.Host)
		if err == nil {
			f.forwardAll[name] = []interface{}{listen.Scheme, listener, listen, router}
			go f.loopForward(listener, name, listen, router)
			InfoLog("Forward(%v) start tcp forward on %v success by %v->%v", f.Name, listener.Addr(), listen, router)
		}
	}
	return
}

// StopForward will forward address to uri
func (f *Forward) StopForward(name string) (err error) {
	InfoLog("Forward(%v) stop forward by name:%v", f.Name, name)
	f.forwardLck.Lock()
	forward := f.forwardAll[name]
	delete(f.forwardAll, name)
	f.forwardLck.Unlock()
	if len(forward) > 0 {
		err = forward[1].(io.Closer).Close()
	}
	return
}

func (f *Forward) loopForward(l net.Listener, name string, listen *url.URL, uri string) {
	defer func() {
		f.forwardLck.Lock()
		delete(f.forwardAll, name)
		f.forwardLck.Unlock()
	}()
	var err error
	var piper xio.Piper
	var conn net.Conn
	InfoLog("Forward(%v) proxy forward(%v->%v) accept runner is starting", f.Name, l.Addr(), uri)
	for {
		conn, err = l.Accept()
		if err != nil {
			break
		}
		DebugLog("Forward(%v) accepting forward(%v->%v) connection from %v", f.Name, l.Addr(), uri, conn.RemoteAddr())
		piper, err = f.Dialer.DialPiper(uri, f.BufferSize)
		if err == nil {
			DebugLog("Forward(%v) proxy forward(%v->%v) success", f.Name, l.Addr(), uri)
			go f.procForward(l, name, piper, conn, uri)
		} else {
			WarnLog("Forward(%v) proxy forward(%v->%v) fail with %v", f.Name, l.Addr(), uri, err)
			conn.Close()
		}
	}
	l.Close()
	InfoLog("Forward(%v) proxy forward(%v->%v) accept runner is stopped", f.Name, l.Addr(), uri)
}

func (f *Forward) procForward(l net.Listener, name string, piper xio.Piper, conn net.Conn, uri string) {
	defer func() {
		if perr := recover(); perr != nil {
			ErrorLog("Forward(%v) transfer forward (%v->%v) is pance with %v, callstack is \n%v", f.Name, l.Addr(), uri, perr, xdebug.CallStack())
		}
	}()
	piper.PipeConn(conn, uri)
}

// Stop will stop all
func (f *Forward) Stop() (err error) {
	InfoLog("Forward(%v) is closing", f.Name)
	f.forwardLck.RLock()
	for key, forward := range f.forwardAll {
		forward[1].(net.Listener).Close()
		InfoLog("Forward(%v) forwad %v is closed", f.Name, key)
	}
	f.forwardAll = map[string][]interface{}{}
	f.forwardLck.RUnlock()
	return
}
