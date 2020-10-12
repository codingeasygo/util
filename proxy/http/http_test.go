package http

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"strings"
	"testing"
	"time"
)

func init() {
	SetLogLevel(LogLevelDebug)
	go http.ListenAndServe(":6060", nil)
}

func TestProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "abc")
	}))
	server := NewServer()
	_, err := server.Start(":8011")
	if err != nil {
		t.Error(err)
		return
	}
	{ //CONNECT
		client := http.Client{
			Transport: &http.Transport{
				Dial: func(network, address string) (conn net.Conn, err error) {
					conn, err = Dial(":8011", address)
					return
				},
			},
		}
		{ //ok
			resp, err := client.Get(ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if string(data) != "abc" {
				t.Error(string(data))
				return
			}
		}
		{ //error
			_, err := client.Get("http://127.0.0.1:233")
			if err == nil {
				t.Error(err)
				return
			}
		}
	}
	{ //NORMAL
		client := http.Client{
			Transport: &http.Transport{
				Dial: func(network, address string) (conn net.Conn, err error) {
					conn, err = net.Dial("tcp", ":8011")
					return
				},
			},
		}
		{ //ok
			req, _ := http.NewRequest("GET", ts.URL, nil)
			req.Header.Add("Proxy-Connection", "keep-alive")
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if string(data) != "abc" {
				t.Error(string(data))
				return
			}
		}
		{ //not port
			req, _ := http.NewRequest("GET", "http://www.bing.com", nil)
			req.Header.Add("Proxy-Connection", "keep-alive")
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if !strings.Contains(string(data), "bing.com") {
				t.Error(string(data))
				return
			}
		}
		{ //error
			req, _ := http.NewRequest("GET", "http://127.0.0.1:233", nil)
			req.Header.Add("Proxy-Connection", "keep-alive")
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode == 200 {
				t.Error(err)
				return
			}
		}
	}
	{ //info
		client := http.Client{}
		{ //ok
			req, _ := http.NewRequest(http.MethodHead, "http://127.0.0.1:8011", nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				t.Error("error")
				return
			}
			resp.Body.Close()
		}
		{ //500
			req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8011", nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			if resp.StatusCode != http.StatusInternalServerError {
				t.Error("error")
				return
			}
			resp.Body.Close()
		}
	}
	{ //ERROR
		conn, _ := net.Dial("tcp", ":8011")
		time.Sleep(10 * time.Millisecond)
		conn.Close()
		Dial("127.0.0.1:233", "")
		Dial("127.0.0.1:8011", "%2f")
	}
	server.Stop()
	go func() {
		server.Run(":8011")
	}()
	time.Sleep(10 * time.Millisecond)
	server.Stop()
	// conn,err:=Dial(":8011",ts.URL)
}
