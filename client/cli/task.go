package cli

import (
	"fmt"
	"net/http"
	"os"
)

func newShellcode(path string) ([]byte, error) {
	sc, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(sc)
	//checks if the file is a hex file
	if contentType == "text/plain; charset=utf-8" {
		return sc, nil
		// check if the file is a binary
	} else if contentType == "application/octet-stream" {
		hex := fmt.Sprintf("%x ", sc)
		return []byte(hex), nil
	}
	return nil, nil
}
