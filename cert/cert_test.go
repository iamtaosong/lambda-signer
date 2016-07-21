package cert

import (
	"testing"
)

func TestGenerateCert(t *testing.T) {
	IPs := []string{"127.0.0.1"}
	hosts := []string{"example.com"}

	opts := &Options{
		Hosts: append(IPs, hosts...),
		Org:   "test",
	}

	cert := GenerateCert(opts)
	if len(cert.IPAddresses) != 1 {
		t.Fatalf("Wrong number of IPs found: %d instead of 1", len(cert.IPAddresses))
	}

	if string(cert.IPAddresses[0]) != "127.0.0.1" {
		t.Fatalf("Wrong IP found")
	}

	if cert.Subject.Organization != opts.Org {
		t.Fatalf("Wrong organisation found")
	}
}
