package client

import (
	"KittyStager/internal/kitten"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
)

func printKittens(kittens map[string]*kitten.Kitten) error {
	if len(kittens) < 2 {
		return errors.New("No kittens to show")
	}
	fmt.Printf("%s\n\n", color.Green("[*] Kittens:"))
	fmt.Printf("%s\n", color.Green("Name:\tIp:\t\tHostname:\t\tLast seen:\tSleep:\tAlive:"))
	fmt.Printf("%s\n", color.Green("═════\t═══\t\t═════════\t\t══════════\t══════\t══════"))

	for name, k := range kittens {
		var e string
		if name != "" {
			r := k.GetRecon()
			if k.GetAlive() {
				e = "Yes"
				fmt.Printf("%s\t%s\t%s\t\t%s\t%d\t%s\n",
					color.Yellow(name),
					color.Yellow(r.GetIp()),
					color.Yellow(r.GetHostname()),
					color.Yellow(k.GetLastSeen().Format("15:04:05")),
					color.Yellow(k.GetSleep()),
					color.Yellow(e))

			} else {
				e = "No 💀"
				fmt.Printf("%s\t%s\t%s\t\t%s\t%d\t%s\n",
					color.Red(name),
					color.Red(r.GetIp()),
					color.Red(r.GetHostname()),
					color.Red(k.GetLastSeen().Format("15:04:05")),
					color.Red(k.GetSleep()),
					color.Red(e))
			}
		}
	}
	fmt.Println()
	return nil
}

func printKittenInfo(kitten kitten.Kitten) {
	r := kitten.Recon
	fmt.Printf("%s\n\n", color.Green("[*] Kitten:"))
	fmt.Println("{")
	fmt.Printf("\t\"%s\": \"%s\",\n", color.Green("Name"), color.Yellow(kitten.Name))
	fmt.Printf("\t\"%s\": \"%d\",\n", color.Green("Sleep"), color.Yellow(kitten.Sleep))
	fmt.Printf("\t\"%s\": \"%s\",\n", color.Green("LastSeen"), color.Yellow(kitten.LastSeen.Format("15:04:05")))
	fmt.Printf("\t\"%s\": \"%t\",\n", color.Green("Alive"), color.Yellow(kitten.GetAlive()))
	fmt.Printf("\t\"%s\": \"%s\",\n", color.Green("Key"), color.Yellow(kitten.Key))
	fmt.Printf("\t\"%s\": [\n", color.Green("Tasks"))
	for _, b := range kitten.Tasks {
		fmt.Println("\t\t\t{")
		fmt.Printf("\t\t\t\"%s\": \"%s\",\n", color.Green("Tag"), color.Yellow(b.Tag))
		if len(b.Payload) > 10 {
			fmt.Printf("\t\t\t\"%s\": \"%s...\",\n", color.Green("Payload"), color.Yellow(b.Payload[:10]))
		} else {
			fmt.Printf("\t\t\t\"%s\": \"%s\",\n", color.Green("Payload"), color.Yellow(string(b.Payload)))
		}
		fmt.Println("\t\t\t}")
	}
	fmt.Println("\t\t]\n")
	fmt.Printf("\t\"%s\": {\n", color.Green("Recon"))
	fmt.Printf("\t\t\"%s\": \"%s\",\n", color.Green("hostname"), color.Yellow(r.Hostname))
	fmt.Printf("\t\t\"%s\": \"%s\",\n", color.Green("username"), color.Yellow(r.Username))
	fmt.Printf("\t\t\"%s\": \"%s\",\n", color.Green("domain"), color.Yellow(r.Domain))
	fmt.Printf("\t\t\"%s\": \"%s\",\n", color.Green("ip"), color.Yellow(r.Ip))
	fmt.Printf("\t\t\"%s\": \"%d\",\n", color.Green("pid"), color.Yellow(r.Pid))
	fmt.Printf("\t\t\"%s\": \"%s\",\n", color.Green("pname"), color.Yellow(r.PName))
	fmt.Printf("\t\t\"%s\": \"%s\",\n", color.Green("path"), color.Yellow(r.Path))
	fmt.Println("\t\t}")
	fmt.Println("}")
}
