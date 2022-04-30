package xnet

import (
	"fmt"
	"testing"
)

func TestWebsocketDialer(t *testing.T) {
	dialer := NewWebsocketDialer()
	conn, err := dialer.Dial("wss://wmservice:YiGWXa@v100.wmservice.sxbastudio.com/_s/docker/logs")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(conn)
}
