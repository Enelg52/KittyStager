package cli

import (
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"KittyStager/internal/task/ps"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"sort"
	"strings"
)

// https://github.com/ZephrFish/edr-checker/blob/master/Invoke-EDRChecker.ps1
var (
	av = []string{"activeconsole",
		"amsi.dll",
		"authtap",
		"avast",
		"avecto",
		"canary",
		"carbon",
		"cb.exe",
		"ciscoamp",
		"cisco amp",
		"countertack",
		"cramtray",
		"crssvc",
		"crowdstrike",
		"csagent",
		"csfalcon",
		"csshell",
		"cybereason",
		"cyclorama",
		"cylance",
		"cyoptics",
		"cyupdate",
		"cyvera",
		"cyserver",
		"cytray",
		"defendpoint",
		"defender",
		"eectrl",
		"emcoreservice",
		"emsystem",
		"endgame",
		"fireeye",
		"forescout",
		"fortiedr",
		"groundling",
		"GRRservice",
		"healthservice",
		"inspector",
		"ivanti",
		"kaspersky",
		"lacuna",
		"logrhythm",
		"logcollector",
		"malware",
		"mandiant",
		"mcafee",
		"monitoringhost",
		"morphisec",
		"mpcmdrun",
		"msascuil",
		"msmpeng",
		"mssense",
		"msmpeng",
		"nissrv",
		"ntrtscan",
		"osquery",
		"Palo Alto Networks",
		"pgeposervice",
		"pgsystemtray",
		"privilegeguard",
		"procwall",
		"protectorservice",
		"qradar",
		"redcloak",
		"secureconnector",
		"secureworks",
		"securityhealthservice",
		"semlaunchsvc",
		"senseir",
		"sense",
		"sentinel",
		"sepliveupdate",
		"sisidsservice",
		"sisipsservice",
		"sisipsutil",
		"smc.exe",
		"smcgui",
		"snac64",
		"sophos",
		"splunk",
		"srtsp",
		"symantec",
		"symcorpui",
		"symefasi",
		"sysinternal",
		"sysmon",
		"tanium",
		"tda.exe",
		"tdawork",
		"tmlisten",
		"tmbmsrv",
		"tmssclient",
		"tmccsf",
		"tpython",
		"trend",
		"watchdogagent",
		"wincollect",
		"windowssensor",
		"wireshark",
		"xagt"}
)

func printKittens(kittens map[string]*kitten.Kitten) error {
	if len(kittens) < 2 {
		return errors.New("No kittens to show")
	}
	fmt.Printf("%s\n\n", color.BrightGreen("[*] Kittens:"))
	fmt.Printf("%s\n", color.BrightGreen("Name:\tIp:\t\tHostname:\t\tLast seen:\tSleep:\tAlive:"))
	fmt.Printf("%s\n", color.BrightGreen("═════\t═══\t\t═════════\t\t══════════\t══════\t══════"))

	for name, k := range kittens {
		var e string
		if name != "" {
			r := k.GetRecon()
			if k.GetAlive() {
				e = "Yes"
				fmt.Printf("%s\t%s\t%s\t\t%s\t%d\t%s\n",
					name,
					r.GetIp(),
					r.GetHostname(),
					k.GetLastSeen().Format("15:04:05"),
					k.GetSleep(),
					e,
				)

			} else {
				e = "No 💀"
				fmt.Printf("%s\t%s\t%s\t\t%s\t%d\t%s\n",
					color.BrightRed(name),
					color.BrightRed(r.GetIp()),
					color.BrightRed(r.GetHostname()),
					color.BrightRed(k.GetLastSeen().Format("15:04:05")),
					color.BrightRed(k.GetSleep()),
					color.BrightRed(e))
			}
		}
	}
	fmt.Println()
	return nil
}

func printKittenInfo(kitten kitten.Kitten) {
	r := kitten.Recon
	fmt.Printf("%s\n\n", color.BrightGreen("[*] Kitten:"))
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Name"), kitten.Name)
	fmt.Printf("%s:\t\t%d\n", color.BrightGreen("Sleep"), kitten.Sleep)
	fmt.Printf("%s:\t%s\n", color.BrightGreen("LastSeen"), kitten.LastSeen.Format("15:04:05"))
	fmt.Printf("%s:\t\t%v\n", color.BrightGreen("Alive"), kitten.GetAlive())
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Key"), kitten.Key)
	fmt.Printf("%s:\t%s\n", color.BrightGreen("Hostname"), r.Hostname)
	fmt.Printf("%s:\t%s\n", color.BrightGreen("Username"), r.Username)
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Domain"), r.Domain)
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Ip"), r.Ip)
	fmt.Printf("%s:\t\t%d\n", color.BrightGreen("Pid"), r.Pid)
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Pname"), r.PName)
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Path"), r.Path)

}

func printPS(t *task.Task, pid int) error {
	prlist := ps.NewProcessList(nil)
	err := prlist.UnmarshallProcessList(t.Payload)
	if err != nil {
		return err
	}
	sort.Slice(prlist.Process, func(p, q int) bool {
		return prlist.Process[p].Pid < prlist.Process[q].Pid
	})
	fmt.Printf("\n%5s\t%5s\t%s\n", color.BrightGreen("Ppid:"), color.BrightGreen("Pid:"), color.BrightGreen("Name:"))
	fmt.Printf("%s\t%s\t%s\n", color.BrightGreen("═════"), color.BrightGreen("═════"), color.BrightGreen("═════"))
	for _, p := range prlist.Process {
		//highlight the current process in green
		if p.Pid == pid {
			fmt.Printf("%5d\t%5d\t%s\n", color.BrightGreen(p.Ppid), color.BrightGreen(p.Pid), color.BrightGreen(p.Name))
		} else if contains(av, p.Name) {
			fmt.Printf("%5d\t%5d\t%s\n", color.BrightRed(p.Ppid), color.BrightRed(p.Pid), color.BrightRed(p.Name))
		} else {
			fmt.Printf("%5d\t%5d\t%s\n", p.Ppid, p.Pid, p.Name)
		}
	}
	return nil
}

func printAV(t *task.Task) {
	fmt.Printf("\n%s\n\n", color.BrightGreen("[*] AV/EDR:"))
	fmt.Printf("%s\n", string(t.Payload))
}

func printTasks(t []*task.Task) {
	fmt.Printf("%s\n\n", color.BrightGreen("[*] Tasks:"))
	fmt.Printf("%s\n", color.BrightGreen("ID:\tTag:\tPayload:"))
	fmt.Printf("%s\n", color.BrightGreen("═══\t════\t════════"))

	for i, b := range t {
		fmt.Printf("%2d\t%s\t", i, b.Tag)
		if len(b.Payload) > 10 {
			fmt.Printf("%s...\n", b.Payload[:10])
		} else {
			fmt.Printf("%s\n", string(b.Payload))
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(strings.ToLower(str), v) {
			return true
		}
	}

	return false
}
