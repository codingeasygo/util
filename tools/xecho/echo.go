package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type ll struct {
	net.Listener
}

func (l *ll) Accept() (conn net.Conn, err error) {
	conn, err = l.Listener.Accept()
	if err == nil {
		fmt.Printf("ACCEPT:%v\n", conn.RemoteAddr())
	}
	return
}

func main() {
	http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		fmt.Printf("connection from %v is starting\n", ws.RemoteAddr())
		buffer := make([]byte, 32*1024)
		for {
			n, err := ws.Read(buffer)
			if err != nil {
				break
			}
			fmt.Printf("RECV:%v\n", (buffer[0:n]))
			_, err = ws.Write(buffer[0:n])
			if err != nil {
				break
			}
			// binary.BigEndian.PutUint32(buffer, 64)
			// buffer[0] = 3
			// _, err = ws.Write(buffer[:4])
			// if err != nil {
			// 	break
			// }
			// for i := 0; i < 6; i++ {
			// 	time.Sleep(100 * time.Millisecond)
			// 	copy(buffer, []byte("0123456789"))
			// 	_, err = ws.Write(buffer[:10])
			// 	if err != nil {
			// 		break
			// 	}
			// }
		}
		fmt.Printf("connection from %v is done\n", ws.RemoteAddr())
	}))
	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	server := &http.Server{Addr: os.Args[1], Handler: http.DefaultServeMux}
	fmt.Println(server.Serve(&ll{Listener: ln}))
}
