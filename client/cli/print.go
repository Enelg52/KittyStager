package cli

import (
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"KittyStager/internal/task/ps"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"sort"
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
	fmt.Printf("%s:\t\t%s\n", color.Green("Name"), kitten.Name)
	fmt.Printf("%s:\t\t%d\n", color.Green("Sleep"), kitten.Sleep)
	fmt.Printf("%s:\t%s\n", color.Green("LastSeen"), kitten.LastSeen.Format("15:04:05"))
	fmt.Printf("%s:\t\t%v\"\n", color.Green("Alive"), kitten.GetAlive())
	fmt.Printf("%s:\t\t%s\n", color.Green("Key"), kitten.Key)
	fmt.Printf("%s:\t%s\n", color.Green("Hostname"), r.Hostname)
	fmt.Printf("%s:\t%s\n", color.Green("Username"), r.Username)
	fmt.Printf("%s:\t\t%s\n", color.Green("Domain"), r.Domain)
	fmt.Printf("%s:\t\t%s\n", color.Green("Ip"), r.Ip)
	fmt.Printf("%s:\t\t%d\n", color.Green("Pid"), r.Pid)
	fmt.Printf("%s:\t\t%s\n", color.Green("Pname"), r.PName)
	fmt.Printf("%s:\t\t%s\n", color.Green("Path"), r.Path)

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
	fmt.Printf("\n%5s\t%5s\t%s\n", color.Green("Ppid:"), color.Green("Pid:"), color.Green("Name:"))
	fmt.Printf("%s\t%s\t%s\n", color.Green("═════"), color.Green("═════"), color.Green("═════"))
	for _, p := range prlist.Process {
		//highlight the current process in green
		if p.Pid == pid {
			fmt.Printf("%5d\t%5d\t%s\n", color.Green(p.Ppid), color.Green(p.Pid), color.Green(p.Name))
		} else {
			fmt.Printf("%5d\t%5d\t%s\n", p.Ppid, p.Pid, p.Name)
		}
	}
	return nil
}

func printAV(t *task.Task) {
	fmt.Printf("\n%s\n\n", color.Green("[*] AV/EDR:"))
	fmt.Printf("%s\n", string(t.Payload))
}

func printTasks(t []*task.Task) {
	fmt.Printf("%s\n\n", color.Green("[*] Tasks:"))
	fmt.Printf("%s\n", color.Green("ID:\tTag:\tPayload:"))
	fmt.Printf("%s\n", color.Green("═══\t════\t════════"))

	for i, b := range t {
		fmt.Printf("%2d\t%s\t", i, b.Tag)
		if len(b.Payload) > 10 {
			fmt.Printf("%s...\n", b.Payload[:10])
		} else {
			fmt.Printf("%s\n", string(b.Payload))
		}
	}
}
