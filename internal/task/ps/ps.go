package ps

import "encoding/json"

type ProcessList struct {
	process []Process
}

type Process struct {
	ppid int
	pid  int
	name string
}

func NewProcessList(p *[]Process) *ProcessList {
	return &ProcessList{process: *p}
}

func NewProcess(ppid, pid int, name string) *Process {
	return &Process{
		ppid: ppid,
		pid:  pid,
		name: name,
	}
}

func (processList *ProcessList) MarshallTask() ([]byte, error) {
	t, err := json.Marshal(processList)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (processList *ProcessList) UnmarshallTask(j []byte) error {
	err := json.Unmarshal(j, &processList)
	if err != nil {
		return err
	}
	return nil
}
