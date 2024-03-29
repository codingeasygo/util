package xcrypto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

func init() {
	// exec.Command("bash", "-c", "./openssl.sh").Output()
}

func testWebCert(t *testing.T, caPEM []byte, serverCert, clientCert tls.Certificate) {
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(caPEM))
	if !ok {
		panic("failed to parse ca certificate")
	}
	var server *http.Server
	{
		server = &http.Server{}
		server.Addr = ":8122"
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		}
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "ok")
		})
		go server.ListenAndServeTLS("", "")
		defer server.Close()
		time.Sleep(10 * time.Millisecond)
	}
	{
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      certPool,
			},
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial(network, server.Addr)
			},
		}
		client := &http.Client{Transport: transport}
		resp, err := client.Get("https://a.test.com")
		if err != nil {
			t.Error(err)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if string(data) != "ok" {
			t.Errorf("%v", string(data))
			return
		}
		fmt.Println("-->", string(data))
	}
	{
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      certPool,
			},
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial(network, server.Addr)
			},
		}
		client := &http.Client{Transport: transport}
		resp, err := client.Get("https://127.0.0.1")
		if err != nil {
			t.Error(err)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if string(data) != "ok" {
			t.Errorf("%v", string(data))
			return
		}
		fmt.Println("-->", string(data))
	}
}

// func TestGenerateWeb(t *testing.T) {
// 	cert, certPEM, _, err := GenerateWeb(nil, nil, "a.test.com", "127.0.0.1", 2048)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	testWebCert(t, cert, certPEM)
// }

func TestGenerateRoot(t *testing.T) {
	rootCert, rootPriv, rootCertPEM, _, err := GenerateRootCA([]string{"test"}, "Root CA", 2048)
	if err != nil {
		t.Error(err)
		return
	}
	serverCert, _, _, err := GenerateWeb(rootCert, rootPriv, false, "test", "a.test.com", "127.0.0.1", 2048)
	if err != nil {
		t.Error(err)
		return
	}
	clientCert, _, _, err := GenerateWeb(rootCert, rootPriv, true, "test", "a.test.com", "127.0.0.1", 2048)
	if err != nil {
		t.Error(err)
		return
	}
	testWebCert(t, rootCertPEM, serverCert, clientCert)
}

// func TestOpensslRoot(t *testing.T) {
// 	certPEM, err := ioutil.ReadFile("ca.pem")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	testWebCert(t, cert, certPEM)
// }

func TestGenerateWebServerClient(t *testing.T) {
	_, _, _, _, _, _, _, _, _, _, err := GenerateWebServerClient("test", "ca", "test", "", 2048)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLoadX509KeyPair(t *testing.T) {
	_, _, err := LoadX509KeyPair("ca.pem", "ca.key")
	if err != nil {
		t.Error(err)
		return
	}
	_, _, err = LoadX509KeyPair("ca.pem", "ca.keyxx")
	if err == nil {
		t.Error(err)
		return
	}
	_, _, err = LoadX509KeyPair("ca.pemxx", "ca.key")
	if err == nil {
		t.Error(err)
		return
	}
}
