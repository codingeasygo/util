package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/codingeasygo/util/xnet"
)

var waiter = sync.WaitGroup{}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf(`Usage: transport <local> <remote> <local> <remote> ...`)
		return
	}
	mappings := os.Args[1:]
	n := len(mappings) / 2
	for i := 0; i < n; i++ {
		waiter.Add(1)
		go mapping(mappings[i*2], mappings[i*2+1])
	}
	waiter.Done()
}

func mapping(local, remote string) {
	defer waiter.Done()
	var transporter xnet.Transporter
	if strings.HasPrefix(remote, "tcp://") {
		transporter = xnet.RawDialerF(net.Dial)
	} else if strings.HasPrefix(remote, "ws://") || strings.HasPrefix(remote, "wss://") {
		transporter = xnet.NewWebsocketDialer()
	} else {
		err := fmt.Errorf("not supported remote %v", remote)
		ErrorLog("mapping %v to %v fail with %v", local, remote, err)
		return
	}
	ln, err := net.Listen("tcp", local)
	if err != nil {
		ErrorLog("mapping %v to %v fail with %v", local, remote, err)
		return
	}
	InfoLog("start mapping %v to %v", local, remote)
	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			break
		}
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			transporter.Transport(conn, remote)
		}()
	}
	InfoLog("mapping %v to %v is done with %v", local, remote, err)
}
