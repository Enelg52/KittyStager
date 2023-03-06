package cli

import (
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"KittyStager/internal/task/priv"
	"KittyStager/internal/task/ps"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"sort"
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
					color.BrightWhite(name),
					color.BrightWhite(r.GetIp()),
					color.BrightWhite(r.GetHostname()),
					color.BrightWhite(k.GetLastSeen().Format("15:04:05")),
					color.BrightWhite(k.GetSleep()),
					color.BrightWhite(e),
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
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Name"), color.BrightWhite(kitten.Name))
	fmt.Printf("%s:\t\t%d\n", color.BrightGreen("Sleep"), color.BrightWhite(kitten.Sleep))
	fmt.Printf("%s:\t%s\n", color.BrightGreen("LastSeen"), color.BrightWhite(kitten.LastSeen.Format("15:04:05")))
	fmt.Printf("%s:\t\t%v\n", color.BrightGreen("Alive"), color.BrightWhite(kitten.GetAlive()))
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Key"), color.BrightWhite(kitten.Key))
	fmt.Printf("%s:\t%s\n", color.BrightGreen("Hostname"), color.BrightWhite(r.Hostname))
	fmt.Printf("%s:\t%s\n", color.BrightGreen("Username"), color.BrightWhite(r.Username))
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Domain"), color.BrightWhite(r.Domain))
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Ip"), color.BrightWhite(r.Ip))
	fmt.Printf("%s:\t\t%d\n", color.BrightGreen("Pid"), color.BrightWhite(r.Pid))
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Pname"), color.BrightWhite(r.PName))
	fmt.Printf("%s:\t\t%s\n", color.BrightGreen("Path"), color.BrightWhite(r.Path))

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
			fmt.Printf("%5d\t%5d\t%s\n", color.BrightWhite(p.Ppid), color.BrightWhite(p.Pid), color.BrightWhite(p.Name))
		}
	}
	return nil
}

func printAV(t *task.Task) {
	fmt.Printf("\n%s\n\n", color.BrightGreen("[*] AV/EDR:"))
	fmt.Printf("%s\n", color.BrightWhite(string(t.Payload)))
}

func printTasks(t []*task.Task) {
	fmt.Printf("%s\n\n", color.BrightGreen("[*] Tasks:"))
	fmt.Printf("%s\n", color.BrightGreen("ID:\tTag:\tPayload:"))
	fmt.Printf("%s\n", color.BrightGreen("═══\t════\t════════"))

	for i, b := range t {
		fmt.Printf("%2d\t%s\t", color.BrightWhite(i), color.BrightWhite(b.Tag))
		if len(b.Payload) > 10 {
			fmt.Printf("%s...\n", color.BrightWhite(b.Payload[:10]))
		} else {
			fmt.Printf("%s\n", color.BrightWhite(string(b.Payload)))
		}
	}
}

func printPriv(t *task.Task) {
	p := priv.NewPrivileges(nil, "")
	err := p.UnmarshallPrivileges(t.Payload)
	if err != nil {
		fmt.Println("[!] Error", err)
		return
	}
	fmt.Printf("\n%s\n\n", color.BrightGreen("[*] Privileges:"))
	fmt.Printf("%s\n", color.BrightGreen("Privileges:"))
	fmt.Printf("%s\n", color.BrightGreen("═══════════"))
	for _, pr := range p.Priv {
		fmt.Printf("%-29s\t%#v\t%s\n", color.BrightWhite(pr.GetName()), color.BrightWhite(pr.GetEnable()), color.BrightWhite(pr.GetDescription()))
	}
	fmt.Printf("\n%s\n", color.BrightGreen("Integrity:"))
	fmt.Printf("%s\n", color.BrightGreen("══════════"))
	fmt.Println(color.BrightWhite(p.Integrity))
}

func printHelpInt() {
	fmt.Printf("%s\n\n", color.BrightGreen("Help :"))
	fmt.Printf("\n%5s\t%5s\n", color.BrightGreen("Command:"), color.BrightGreen("Description:"))
	fmt.Printf("%s\t%s\n", color.BrightGreen("════════"), color.BrightGreen("════════════"))
	fmt.Printf("%s\n", color.BrightWhite("back\t\tGo back to the main menu\n"+
		"help\t\tPrint the help menu\n"+
		"tasks\t\tGet all the current tasks for the current kitten\n"+
		"shellcode\tInject shellcode in new process\n"+
		"sleep\t\tSet sleep time\n"+
		"ps\t\tGet process list\n"+
		"av\t\tGet AV/EDR with wmi\n"+
		"priv\t\tGet privileges and integrity level\n"+
		"info\t\tShow all the kitten info\n"+
		"kill\t\tKill the kitten :(\n"+
		"exit\t\tExit the client"))
}

func printHelpMain() {
	fmt.Printf("%s\n\n", color.BrightGreen("Help :"))
	fmt.Printf("\n%5s\t%5s\n", color.BrightGreen("Command:"), color.BrightGreen("Description:"))
	fmt.Printf("%s\t%s\n", color.BrightGreen("════════"), color.BrightGreen("════════════"))
	fmt.Printf("%s\n", color.BrightWhite("exit\t\tExit the client\n"+
		"help\t\tPrint the help menu\n"+
		"config\t\tShow the server config\n"+
		"logs\t\tDisplay the log in real time\n"+
		"kittens\t\tShow all kittens\n"+
		"interact\tInteract with a kitten\n"))
}
