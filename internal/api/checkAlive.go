package api

import (
	"time"
)

func CheckAlive(name string) {
	k := Kittens[name]

	for {
		time.Sleep(1 * time.Second)
		if k.GetAlive() {
			t := time.Since(k.GetLastSeen())
			sleepTime := time.Duration(k.GetSleep()) * time.Second
			if t > sleepTime+5*time.Second {
				k.SetAlive(false)
			}
		}
	}
}
