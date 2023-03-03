package ps

import "encoding/json"

type ProcessList struct {
	Process []Process `json:"process"`
}

type Process struct {
	Ppid int    `json:"Ppid"`
	Pid  int    `json:"Pid"`
	Name string `json:"Name"`
}

func NewProcessList(p []Process) *ProcessList {
	return &ProcessList{Process: p}
}

func NewProcess(ppid, pid int, name string) *Process {
	return &Process{
		Ppid: ppid,
		Pid:  pid,
		Name: name,
	}
}

func (processList *ProcessList) MarshallProcessList() ([]byte, error) {
	t, err := json.Marshal(processList)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (processList *ProcessList) UnmarshallProcessList(j []byte) error {
	err := json.Unmarshal(j, &processList)
	if err != nil {
		return err
	}
	return nil
}
