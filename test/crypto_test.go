package test

import (
	"KittyStager/internal/crypto"
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
