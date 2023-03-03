package test

import (
	"KittyStager/internal/task/recon"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	Hostname string
	Username string
	Domain   string
	Ip       string
	//process
	Pid   int
	PName string
	Path  string
)

func reconBeforeAll() {
	Hostname = "host"
	Username = "user"
	Domain = "domain"
	Ip = "ip"
	//process
	Pid = 1
	PName = "pName"
	Path = "path"
}

func TestMarshallUnmashallRecon(t *testing.T) {
	t.Parallel()
	// given
	recon := recon.NewRecon(Hostname, Username, Domain, Ip, PName, Path, Pid)
	// when
	b, err := recon.MarshallTask()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = recon.UnmarshallTask(b)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	if recon.Hostname != Hostname {
		t.Errorf("Failed")
	}
}

func TestEncryptDecryptRecon(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(Hostname, Username, Domain, Ip, PName, Path, Pid)
	// when
	e, err := r.EncryptTask(Key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = r.DecryptTask(e, Key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, r.Hostname, Hostname)
}
