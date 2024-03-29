package xcrypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"
)

func generateRSA(isCA bool, template, parent *x509.Certificate, rootKey *rsa.PrivateKey, commonName string, dnsNames []string, ipAddresses []net.IP, bits int) (cert *x509.Certificate, privKey *rsa.PrivateKey, certPEM, privPEM []byte, err error) {
	privKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}

	if parent == nil {
		parent = template
	}
	if rootKey == nil {
		rootKey = privKey
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, &privKey.PublicKey, rootKey)
	if err != nil {
		return
	}
	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return
	}

	certBuffer := &bytes.Buffer{}
	certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	pem.Encode(certBuffer, certBlock)
	certPEM = certBuffer.Bytes()

	privBuffer := &bytes.Buffer{}
	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)}
	pem.Encode(privBuffer, privBlock)
	privPEM = privBuffer.Bytes()
	return
}

func GenerateRootCA(org []string, name string, bits int) (cert *x509.Certificate, privKey *rsa.PrivateKey, certPEM, privPEM []byte, err error) {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: org,
			CommonName:   name,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
	}
	cert, privKey, certPEM, privPEM, err = generateRSA(true, template, nil, nil, name, nil, nil, bits)
	return
}

// func GenerateDistCA(parent *x509.Certificate, rootKey *rsa.PrivateKey, name string, bits int) (cert *x509.Certificate, privKey *rsa.PrivateKey, certPEM, privPEM []byte, err error) {
// 	cert, privKey, certPEM, privPEM, err = generateRSA(true, parent, rootKey, name, nil, nil, bits)
// 	return
// }

func GenerateCert(parent *x509.Certificate, rootKey *rsa.PrivateKey, ext []x509.ExtKeyUsage, org []string, commonName string, dnsNames []string, ipAddresses []net.IP, bits int) (cert *x509.Certificate, privKey *rsa.PrivateKey, certPEM, privPEM []byte, err error) {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: org,
			CommonName:   commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           ext,
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
	}
	cert, privKey, certPEM, privPEM, err = generateRSA(false, template, parent, rootKey, commonName, dnsNames, ipAddresses, bits)
	return
}

// GenerateWeb will generate web cert
func GenerateWeb(parent *x509.Certificate, rootKey *rsa.PrivateKey, client bool, org, domain, ip string, bits int) (cert tls.Certificate, certPEM, privPEM []byte, err error) {
	domains := strings.Split(domain, ",")
	ipAddress := []net.IP{}
	if len(ip) > 0 {
		ips := strings.Split(ip, ",")
		for _, v := range ips {
			ipAddress = append(ipAddress, net.ParseIP(v))
		}
	}
	keyUsage := []x509.ExtKeyUsage{}
	if client {
		keyUsage = append(keyUsage, x509.ExtKeyUsageClientAuth)
	} else {
		keyUsage = append(keyUsage, x509.ExtKeyUsageServerAuth)
	}
	_, _, certPEM, privPEM, err = GenerateCert(parent, rootKey, keyUsage, []string{org}, domain, domains, ipAddress, bits)
	if err == nil {
		cert, err = tls.X509KeyPair(certPEM, privPEM)
	}
	return
}

func GenerateWebServerClient(org string, ca, domain, ip string, bits int) (
	rootCert *x509.Certificate, rootKey *rsa.PrivateKey, rootCertPEM, rootKeyPEM []byte,
	server tls.Certificate, serverCertPEM, serverKeyPEM []byte,
	client tls.Certificate, clientCertPEM, clientKeyPEM []byte,
	err error,
) {
	rootCert, rootKey, rootCertPEM, rootKeyPEM, err = GenerateRootCA([]string{org}, ca, bits)
	if err != nil {
		return
	}
	server, serverCertPEM, serverKeyPEM, err = GenerateWeb(rootCert, rootKey, false, org, domain, ip, bits)
	if err != nil {
		return
	}
	client, clientCertPEM, clientKeyPEM, err = GenerateWeb(rootCert, rootKey, true, org, domain, ip, bits)
	return
}

func LoadX509KeyPair(certFile, keyFile string) (cert *x509.Certificate, priv *rsa.PrivateKey, err error) {
	certPEM, err := ioutil.ReadFile(certFile)
	if err != nil {
		return
	}
	keyPEM, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}
	//
	certBlock, _ := pem.Decode(certPEM)
	keyBlock, _ := pem.Decode(keyPEM)
	cert, err = x509.ParseCertificate(certBlock.Bytes)
	if err == nil {
		priv, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	}
	return
}
