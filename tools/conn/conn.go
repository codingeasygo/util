package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: conn uri\n")
		os.Exit(1)
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("parse uri fail with %v\n", err)
		os.Exit(1)
		return
	}
	conn, err := net.Dial(u.Scheme, u.Host)
	if err != nil {
		fmt.Printf("conn uri fail with %v\n", err)
		os.Exit(1)
		return
	}
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}
