// Copyright 2016 Docker, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

type Options struct {
	Hosts        []string
	Org          string
	RawCAKeyPair []byte
	Bits         int
}

type Generator interface {
	GenerateCert(opts *Options) (io.Reader, error)
}

type X509CertGenerator struct{}

func NewX509CertGenerator() Generator {
	return &X509CertGenerator{}
}

func (xcg *X509CertGenerator) newCertificate(org string) (*x509.Certificate, error) {
	now := time.Now()
	// need to set notBefore slightly in the past to account for time
	// skew in the VMs otherwise the certs sometimes are not yet valid
	notBefore := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()-5, 0, 0, time.Local)
	notAfter := notBefore.Add(time.Hour * 24 * 1080)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		BasicConstraintsValid: true,
	}, nil

}

// GenerateCert generates a new certificate signed using the provided
// certificate authority files and stores the result in the certificate
// file and key provided.  The provided host names are set to the
// appropriate certificate fields.
func (xcg *X509CertGenerator) GenerateCert(opts *Options) (io.Reader, error) {
	template, err := xcg.newCertificate(opts.Org)
	if err != nil {
		return nil, err
	}
	// client
	if len(opts.Hosts) == 1 && opts.Hosts[0] == "" {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		template.KeyUsage = x509.KeyUsageDigitalSignature
	} else { // server
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}

		for _, h := range opts.Hosts {
			if ip := net.ParseIP(h); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else {
				template.DNSNames = append(template.DNSNames, h)
			}
		}
	}

	tlsCert, err := tls.X509KeyPair(opts.RawCAKeyPair, opts.RawCAKeyPair)
	if err != nil {
		return nil, err
	}

	priv, err := rsa.GenerateKey(rand.Reader, opts.Bits)
	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, x509Cert, &priv.PublicKey, tlsCert.PrivateKey)
	if err != nil {
		return nil, err
	}

	var (
		certOut bytes.Buffer
	)

	pem.Encode(&certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	pem.Encode(&certOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return &certOut, nil
}
