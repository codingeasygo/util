package http

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "abc")
	}))
	server := NewServer()
	err := server.Start(":8011")
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
				t.Error("error")
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
			resp, err := client.Get(ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if string(data) != "abc" {
				t.Error("error")
				return
			}
		}
		{ //not port
			resp, err := client.Get("http://www.bing.com")
			if err != nil {
				t.Error(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if string(data) == "" {
				t.Error("error")
				return
			}
		}
		{ //error
			resp, err := client.Get("http://127.0.0.1:233")
			if err != nil || resp.StatusCode == 200 {
				t.Error(err)
				return
			}
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
