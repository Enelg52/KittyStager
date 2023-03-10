package test

import (
	"KittyStager/client/cli"
	"KittyStager/internal/config"
	"KittyStager/internal/crypto"
	"KittyStager/kitten/malware"
	"KittyStager/server/api"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	cryptoData        []byte
	cryptoKey         []byte
	cryptoDecryptData []byte
	cryptoUsername    string
	cryptoPassword    string
	cryptoTestKey     string
)

func init() {
	cryptoData = []byte("test msg")
	cryptoDecryptData = []byte("test msg")
	cryptoKey = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJE")
	cryptoUsername = "geko"
	cryptoPassword = "test"
	cryptoTestKey = "testtestte"
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()
	// given
	c := crypto.NewChaCha20()
	// when
	e, err := c.Encrypt(cryptoData, cryptoKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	d, err := c.Decrypt(e, cryptoKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, d, cryptoData)
}

func TestGenerateKey(t *testing.T) {
	t.Parallel()
	key := crypto.GenerateKey([]byte(cryptoTestKey), 32)
	if len(key) != 32 {
		t.Errorf("Invalid confKey lenght")
	}
	key = crypto.GenerateKey([]byte(cryptoTestKey), 4)
	if len(key) != 4 {
		t.Errorf("Invalid confKey lenght")
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()
	// given
	c := crypto.NewChaCha20()
	_, err := c.Decrypt(cryptoDecryptData, cryptoKey)
	if err == nil {
		t.Fail()
	}
}

func TestOpaque(t *testing.T) {
	malConfig := malware.NewConfig("http://127.0.0.1:8080", "getLegit", "postLegit", "reg", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0", cryptoUsername, 1, 0)
	conf, err := config.NewConfig("config.yaml")

	go api.Api(conf)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	err = malware.DoPwreg(cryptoUsername, cryptoPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// opaque auth
	key, err := malware.DoAuth(cryptoUsername, cryptoPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	kitten, err := cli.GetKittens()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, kitten[cryptoUsername].GetKey(), key)

}
