package cryptoUtil

import (
	"errors"
)

// Encrypt cypher the shellcode with chacha20
func Encrypt(payload []byte, key string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("the key needs to be 32 chars long")
	}
	chacha20 := ChaCha20{}
	byteKey := []byte(key)
	cypherPayload, _ := chacha20.Encrypt(payload, byteKey)
	return cypherPayload, nil
}

// Decrypt decode the payload with chacha20
func Decrypt(cypherText []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("the key needs to be 32 chars long")
	}
	chacha20 := ChaCha20{}
	byteKey := key
	plainText, _ := chacha20.Decrypt(cypherText, byteKey)
	return plainText, nil
}

// GenerateKey generate the key with the hostname
func GenerateKey(hostname string, size int) string {
	// generate 32 char key
	key := hostname
	for len(key) < size {
		key = key + hostname
	}
	if len(key) > size {
		key = key[:size]
	}
	return key
}
