package cryptoUtil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt the payload with AES
func Encrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	cypherText := gcm.Seal(nonce, nonce, plainText, nil)
	return cypherText, nil
}

// DecodeAES decode the payload with AES
func DecodeAES(cypherText []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("the key needs to be 32 chars long")
	}
	shellcode, err := decrypt(cypherText, key)
	if err != nil {
		return nil, err
	}
	return shellcode, nil
}

// decrypt the payload with AES
func decrypt(cypherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plainText, err := gcm.Open(nil, cypherText[:gcm.NonceSize()], cypherText[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
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
