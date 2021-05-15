package xnet

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

//Dialer is interface for dial raw connect by string
type Dialer interface {
	Dial(remote string) (raw io.ReadWriteCloser, err error)
}

//WebsocketDialer is an implementation of Dialer by websocket
type WebsocketDialer struct {
	Dialer     RawDialer
	SkipVerify bool
}

//NewWebsocketDialer will create new WebsocketDialer
func NewWebsocketDialer() (dialer *WebsocketDialer) {
	dialer = &WebsocketDialer{
		Dialer: &net.Dialer{},
	}
	return
}

//Dial dial to remote by websocket
func (w WebsocketDialer) Dial(remote string) (raw io.ReadWriteCloser, err error) {
	targetURL, err := url.Parse(remote)
	if err != nil {
		return
	}
	username, password := targetURL.Query().Get("username"), targetURL.Query().Get("password")
	skipVerify := targetURL.Query().Get("skip_verify") == "1" || w.SkipVerify
	timeout, _ := strconv.ParseUint(targetURL.Query().Get("timeout"), 10, 32)
	if timeout < 1 {
		timeout = 5
	}
	var origin string
	if targetURL.Scheme == "wss" {
		origin = fmt.Sprintf("https://%v", targetURL.Host)
	} else {
		origin = fmt.Sprintf("http://%v", targetURL.Host)
	}
	config, err := websocket.NewConfig(targetURL.String(), origin)
	if err == nil {
		if len(username) > 0 && len(password) > 0 {
			config.Header.Set("Authorization", "Basic "+basicAuth(username, password))
		}
		config.TlsConfig = &tls.Config{}
		config.TlsConfig.InsecureSkipVerify = skipVerify
		raw, err = w.dial(config, time.Duration(timeout)*time.Second)
	}
	return
}

var portMap = map[string]string{
	"ws":  "80",
	"wss": "443",
}

func parseAuthority(location *url.URL) string {
	if _, ok := portMap[location.Scheme]; ok {
		if _, _, err := net.SplitHostPort(location.Host); err != nil {
			return net.JoinHostPort(location.Host, portMap[location.Scheme])
		}
	}
	return location.Host
}

func tlsHandshake(rawConn net.Conn, timeout time.Duration, config *tls.Config) (conn *tls.Conn, err error) {
	errChannel := make(chan error, 2)
	time.AfterFunc(timeout, func() {
		errChannel <- fmt.Errorf("timeout")
	})
	conn = tls.Client(rawConn, config)
	go func() {
		errChannel <- conn.Handshake()
	}()
	err = <-errChannel
	return
}

func (w WebsocketDialer) dial(config *websocket.Config, timeout time.Duration) (conn net.Conn, err error) {
	remote := parseAuthority(config.Location)
	rawConn, err := w.Dialer.Dial("tcp", remote)
	if err == nil {
		if config.Location.Scheme == "wss" {
			conn, err = tlsHandshake(rawConn, timeout, config.TlsConfig)
		} else {
			conn = rawConn
		}
		if err == nil {
			conn, err = websocket.NewClient(config, conn)
		}
		if err != nil {
			rawConn.Close()
		}
	}
	return
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
