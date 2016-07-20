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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCACertificate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	// cleanup
	defer os.RemoveAll(tmpDir)

	caCertPath := filepath.Join(tmpDir, "ca.pem")
	caKeyPath := filepath.Join(tmpDir, "key.pem")
	testOrg := "test-org"
	bits := 2048
	gen := NewX509CertGenerator()
	if err := gen.GenerateCACertificate(caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(caCertPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(caKeyPath); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCert(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	// cleanup
	defer os.RemoveAll(tmpDir)

	caCertPath := filepath.Join(tmpDir, "ca.pem")
	caKeyPath := filepath.Join(tmpDir, "key.pem")
	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "cert-key.pem")
	testOrg := "test-org"
	bits := 2048
	gen := NewX509CertGenerator()
	if err := gen.GenerateCACertificate(caCertPath, caKeyPath, testOrg, bits); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(caCertPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(caKeyPath); err != nil {
		t.Fatal(err)
	}

	opts := &Options{
		Hosts:     []string{},
		CertFile:  certPath,
		CAKeyFile: caKeyPath,
		CAFile:    caCertPath,
		KeyFile:   keyPath,
		Org:       testOrg,
		Bits:      bits,
	}

	if err := gen.GenerateCert(opts); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("certificate not created at %s", certPath)
	}

	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("key not created at %s", keyPath)
	}
}
