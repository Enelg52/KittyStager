package test

import (
	"KittyStager/internal/api"
	"KittyStager/internal/config"
	"KittyStager/malware"
	"KittyStager/pkg/crypto"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	data     []byte
	Key      []byte
	username string
	password string
)

func cryptoBeforeAll() {
	data = []byte("test data")
	Key = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJE")
	username = "test"
	password = "test"
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()
	// given
	c := crypto.NewChaCha20()
	// when
	e, err := c.Encrypt(data, Key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	d, err := c.Decrypt(e, Key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, d, data)
}

func TestOpaque(t *testing.T) {
	// given
	c := malware.NewConfig("http://127.0.0.1:8080",
		"getLegit",
		"postLegit",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0",
		"reg",
		"test",
		0,
	)
	conf, err := config.NewConfig("../config.yaml")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	go func() {
		err := api.Api(conf)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
	}()
	// when
	err = malware.DoPwreg(username, password, *c)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	key, err := malware.DoAuth(username, password, *c)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, key, api.Kittens[name].Key)
	assert.Equal(t, name, api.Kittens[name].Name)
}
