package xcrypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

//GenerateRSAPEM will generate rsa by bits
func GenerateRSAPEM(bits int) (cert, priv string, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return
	}
	privPem := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)}
	certPem := &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	out := &bytes.Buffer{}
	//
	out.Reset()
	pem.Encode(out, certPem)
	cert = out.String()
	out.Reset()
	pem.Encode(out, privPem)
	priv = out.String()
	return
}

//GenerateRSA will generate rsa cert
func GenerateRSA(bits int) (cert tls.Certificate, err error) {
	certPEM, privPEM, err := GenerateRSAPEM(bits)
	if err == nil {
		cert, err = tls.X509KeyPair([]byte(certPEM), []byte(privPEM))
	}
	return
}
