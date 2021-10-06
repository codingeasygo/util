package main

import (
	"io/ioutil"
	"os"

	"github.com/codingeasygo/util/xcrypto"
)

func main() {
	_, _, rootCertPEM, rootKeyPEM, _, severCertPEM, serverKeyPEM, _, clientCertPEM, clientKeyPEM, _ := xcrypto.GenerateWebServerClient(os.Args[1], os.Args[2], os.Args[3], 2048)
	ioutil.WriteFile("ca.pem", rootCertPEM, os.ModePerm)
	ioutil.WriteFile("ca.key", rootKeyPEM, os.ModePerm)
	ioutil.WriteFile("server.pem", severCertPEM, os.ModePerm)
	ioutil.WriteFile("server.key", serverKeyPEM, os.ModePerm)
	ioutil.WriteFile("client.pem", clientCertPEM, os.ModePerm)
	ioutil.WriteFile("client.key", clientKeyPEM, os.ModePerm)
}
