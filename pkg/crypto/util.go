package crypto

import "encoding/base64"

func GenerateKey(base []byte, size int) string {
	// generate 32 char key
	key := base
	for len(key) < size {
		key = append(key, base...)
	}
	if len(key) > size {
		key = key[:size]
	}
	keyB64 := base64.StdEncoding.EncodeToString(key)
	return keyB64[:size]
}