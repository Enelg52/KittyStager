package test

import (
	"KittyStager/internal/task/recon"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	reconHostname string
	reconUsername string
	reconDomain   string
	reconIp       string
	//process
	reconPid   int
	reconPName string
	reconPath  string

	reconKey            []byte
	reconWrongKey       []byte
	reconWrongKeyLenght []byte
)

func init() {
	reconHostname = "confHost"
	reconUsername = "user"
	reconDomain = "domain"
	reconIp = "ip"
	//process
	reconPid = 1
	reconPName = "pName"
	reconPath = "path"

	reconKey = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJE")
	reconWrongKey = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJ2")
	reconWrongKeyLenght = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJ2asd")

}

func TestReconSetGet(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(reconHostname, reconUsername, reconDomain, reconIp, reconPName, reconPath, reconPid)
	// when
	assert.Equal(t, r.GetIp(), reconIp)
	assert.Equal(t, r.GetHostname(), reconHostname)
}

func TestMarshallUnmashallRecon(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(reconHostname, reconUsername, reconDomain, reconIp, reconPName, reconPath, reconPid)
	// when
	b, err := r.MarshallRecon()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = r.UnmarshallRecon(b)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	if r.Hostname != reconHostname {
		t.Errorf("Failed")
	}
}

func TestEncryptDecryptRecon(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(reconHostname, reconUsername, reconDomain, reconIp, reconPName, reconPath, reconPid)
	// when
	e, err := r.EncryptTask(reconKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = r.DecryptTask(e, reconKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, r.Hostname, reconHostname)
}

func TestEncryptDecryptReconWrongKey(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(reconHostname, reconUsername, reconDomain, reconIp, reconPName, reconPath, reconPid)
	// when
	e, err := r.EncryptTask(reconKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = r.DecryptTask(e, reconWrongKey)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func TestEncryptReconWrongKeyLenght(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(reconHostname, reconUsername, reconDomain, reconIp, reconPName, reconPath, reconPid)
	// when
	_, err := r.EncryptTask(reconWrongKeyLenght)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func TestDecryptReconWrongKeyLenght(t *testing.T) {
	t.Parallel()
	// given
	r := recon.NewRecon(reconHostname, reconUsername, reconDomain, reconIp, reconPName, reconPath, reconPid)
	// when
	e, err := r.EncryptTask(reconKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = r.DecryptTask(e, reconWrongKeyLenght)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}
