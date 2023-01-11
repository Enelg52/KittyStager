package generate

import (
	"fmt"
	"os/exec"
)

func GoBuild(output, path string) (error, string) {
	cmd := exec.Command("powershell", "go", "build", "-ldflags", "\"-H=windowsgui -s -w\"", "-o", output, path)
	cmdLine := fmt.Sprintf("%s", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err, ""
	}
	return error(nil), cmdLine
}

func GarbleBuild(output, path string) (error, string) {
	cmd1 := exec.Command("garble", "-h")
	err := cmd1.Run()
	if err.Error() != "exit status 2" {
		return err, ""
	}
	cmd2 := exec.Command("garble", "-tiny", "build", "-o", output, path)
	cmdLine := fmt.Sprintf("%s", cmd2.Args)
	fmt.Println(cmdLine)
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Println(cmd2.CombinedOutput())
		return err2, ""
	}
	return error(nil), cmdLine
}

func Signe(signedBinary, path string) (error, string) {
	cmd1 := exec.Command("Mangle", "-h")
	err := cmd1.Run()
	if err != nil {
		return err, ""
	}
	cmd2 := exec.Command("Mangle", "-C", signedBinary, "-M", "-I", path, "-O", path)
	cmdLine := fmt.Sprintf("%s", cmd2.Args)
	err = cmd2.Run()
	if err != nil {
		fmt.Println(cmd2.CombinedOutput())
		return err, ""
	}
	return error(nil), cmdLine
}

func Description(path string) (error, string) {
	cmd1 := exec.Command("goversioninfo")
	err := cmd1.Run()
	if err.Error() != "exit status 1" {
		fmt.Println(err)
		return err, ""
	}
	cmd2 := exec.Command("go", "generate", path)
	cmdLine := fmt.Sprintf("%s", cmd2.Args)
	err = cmd2.Run()
	if err != nil {
		fmt.Println(cmd2.CombinedOutput())
		return err, ""
	}
	return error(nil), cmdLine
}
