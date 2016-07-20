package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"time"
)

// Options holds information about the certificate to be signed
type Options struct {
	Hosts        []string
	Org          string
	RawCAKeyPair []byte
	Bits         int
}

// GenerateX509KeyPair creates a X509 certificate
func GenerateX509KeyPair(opts *Options) (io.Reader, error) {
	cert := GenerateCert(opts)
	key, err := GenerateKey(opts.Bits)
	if err != nil {
		return nil, err
	}

	CAKeyPair, err := tls.X509KeyPair(opts.RawCAKeyPair, opts.RawCAKeyPair)
	if err != nil {
		return nil, err
	}

	derCert, err := Sign(cert, key, CAKeyPair)
	if err != nil {
		return nil, err
	}

	var keyPair bytes.Buffer

	pem.Encode(&keyPair, &pem.Block{Type: "CERTIFICATE", Bytes: derCert})
	pem.Encode(&keyPair, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return &keyPair, nil
}

// GenerateCert creates a new certificate
func GenerateCert(opts *Options) *x509.Certificate {
	now := time.Now()

	notBefore := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)
	notAfter := notBefore.Add(time.Hour * 24 * 1080)

	var (
		IPs   []net.IP
		hosts []string
	)

	for _, h := range opts.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			IPs = append(IPs, []byte(h))
		} else {
			hosts = append(hosts, h)
		}
	}

	return &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			Organization: []string{opts.Org},
		},

		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		BasicConstraintsValid: true,

		IPAddresses: IPs,
		DNSNames:    hosts,

		NotBefore: notBefore,
		NotAfter:  notAfter,
	}
}

// GenerateKey creates the private key
func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// Sign certificate with CA
func Sign(cert *x509.Certificate, key *rsa.PrivateKey, CAKeyPair tls.Certificate) ([]byte, error) {
	caCert, err := x509.ParseCertificate(CAKeyPair.Certificate[0])
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &key.PublicKey, CAKeyPair.PrivateKey)
	if err != nil {
		return nil, err
	}

	return derBytes, nil
}
