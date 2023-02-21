package main

import (
	"KittyStager/internal/task/ps"
	"fmt"
	pl "github.com/mitchellh/go-ps"
)

func main() {
	var processArray []ps.Process
	processList, err := pl.Processes()
	if err != nil {
		return
	}
	// map ages
	for x := range processList {
		var process pl.Process
		process = processList[x]
		p := ps.NewProcess(process.PPid(), process.Pid(), process.Executable())
		processArray = append(processArray, *p)
	}
	list := ps.NewProcessList(&processArray)
	fmt.Println(list)
}
