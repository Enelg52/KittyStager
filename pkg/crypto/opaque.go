package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"github.com/frekui/opaque"
)

var (
	PrivS    *rsa.PrivateKey
	Session1 *opaque.PwRegServerSession
	Session2 *opaque.AuthServerSession
	name     string
	key      string
	Opa      *opaque.User
)

type Req struct {
	Req   string           `json:"Req"`
	Msg1  opaque.PwRegMsg1 `json:"Msg1"`
	Msg3  opaque.PwRegMsg3 `json:"Msg3"`
	AMsg1 opaque.AuthMsg1  `json:"aMsg1"`
	AMsg3 opaque.AuthMsg3  `json:"aMsg3"`
}

func init() {
	PrivS, _ = rsa.GenerateKey(rand.Reader, 512)
}

func HandleAuth(data1 []byte) ([]byte, string, string) {
	var msg1 Req
	if err := json.Unmarshal(data1, &msg1); err != nil {
		return nil, "", ""
	}
	switch msg1.Req {
	case "req1":
		data := req1(msg1)
		return data, "", ""

	case "req2":
		req2(msg1)
		return []byte("ok"), "", ""
	case "auth1":
		data := auth1(msg1)
		return data, "", ""
	case "auth2":
		auth2(msg1)
		return []byte("ok"), name, key
	}
	return nil, "", ""
}

func req1(msg1 Req) []byte {
	var msg2 opaque.PwRegMsg2
	var err error
	Session1, msg2, err = opaque.PwReg1(PrivS, msg1.Msg1)
	if err != nil {
		return nil
	}
	data2, err := json.Marshal(msg2)
	if err != nil {
		return nil
	}
	return data2
}
func req2(msg1 Req) {
	user := opaque.PwReg3(Session1, msg1.Msg3)
	name = user.Username
	Opa = user
}
func auth1(msg1 Req) []byte {
	user := Opa
	var msg2 opaque.AuthMsg2
	var err error
	Session2, msg2, err = opaque.Auth1(PrivS, user, msg1.AMsg1)
	if err != nil {
		return nil
	}
	data2, err := json.Marshal(msg2)
	if err != nil {
		return nil
	}
	return data2
}

func auth2(msg1 Req) {
	sharedSecret, err := opaque.Auth3(Session2, msg1.AMsg3)
	if err != nil {
		return
	}
	key = GenerateKey(sharedSecret, 32)
}
