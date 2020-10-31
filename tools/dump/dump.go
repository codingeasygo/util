package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var protocol, address, out, target, mode string
	flag.StringVar(&protocol, "p", "tcp", "the listen protocol")
	flag.StringVar(&address, "l", "", "the listen address")
	flag.StringVar(&out, "o", "text", "the print out mode")
	flag.StringVar(&target, "t", "", "the target remote address")
	flag.StringVar(&mode, "m", "", "the runner mode")
	flag.Parse()
	if len(address) < 1 {
		flag.PrintDefaults()
		os.Exit(1)
		return
	}
	if protocol == "udp" {
		addr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			fmt.Println("Can't resolve address: ", err)
			os.Exit(1)
		}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		buffer := make([]byte, 32*1024)
		name := "ECHO "
		for {
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				break
			}
			switch out {
			case "text":
				fmt.Printf("%v(%v):\n%v\n\n", name, n, string(buffer[0:n]))
			default:
				fmt.Printf("%v(%v):\n%v\n\n", name, n, buffer[0:n])
			}
			time.Sleep(time.Millisecond)
			_, err = conn.WriteToUDP(buffer[0:n], addr)
			if err != nil {
				break
			}
		}
		return
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		if len(target) > 0 {
			go forwardConn(conn, target, out)
		} else {
			go copyPrint("ECHO ", conn, conn, out)
		}
	}
}

func forwardConn(local net.Conn, target, out string) {
	fmt.Printf("=================\n\nCONN start %v dial to %v\n\n", out, target)
	remote, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("dial to remote %v fail with %v\n", target, err)
		return
	}
	go copyPrint("RECV ", local, remote, out)
	copyPrint("SEND ", remote, local, out)
	fmt.Printf("CONN connection to %v is closed\n\n=================\n\n", target)
}

func copyPrint(name string, dst, src net.Conn, out string) {
	buffer := make([]byte, 32*1024)
	for {
		n, err := src.Read(buffer)
		if err != nil {
			break
		}
		switch out {
		case "text":
			fmt.Printf("%v(%v):\n%v\n\n", name, n, string(buffer[0:n]))
		default:
			fmt.Printf("%v(%v):\n%v\n\n", name, n, buffer[0:n])
		}
		_, err = dst.Write(buffer[0:n])
		if err != nil {
			break
		}
	}
	dst.Close()
}
