package util

import "encoding/json"

// InitialChecks is the struct that contains the initial recon
type InitialChecks struct {
	Hostname   string   `json:"hostname"`
	Username   string   `json:"username"`
	Domain     string   `json:"domain"`
	Ip         string   `json:"ip"`
	KittenName string   `json:"name"`
	Dir        []string `json:"folders,flow"`
	//process
	Pid   int    `json:"pid"`
	PName string `json:"pname"`
	Path  string `json:"path"`
}

func NewInitialChecks() *InitialChecks {
	return &InitialChecks{}
}

// InitUnmarshalJSON unmarshal the json
func InitUnmarshalJSON(j []byte) (*InitialChecks, error) {
	iniCheck := NewInitialChecks()
	err := json.Unmarshal(j, &iniCheck)
	if err != nil {
		return &InitialChecks{}, err
	}
	return iniCheck, nil
}

func (I *InitialChecks) GetHostname() string {
	return I.Hostname
}

func (I *InitialChecks) SetHostname(h string) {
	I.Hostname = h
}

func (I *InitialChecks) GetUsername() string {
	return I.Username
}

func (I *InitialChecks) SetUsername(u string) {
	I.Username = u
}

func (I *InitialChecks) GetDir() []string {
	return I.Dir
}

func (I *InitialChecks) SetDir(d []string) {
	I.Dir = d
}

func (I *InitialChecks) GetIp() string {
	return I.Ip
}

func (I *InitialChecks) SetIp(ip string) {
	I.Ip = ip
}

func (I *InitialChecks) GetKittenName() string {
	return I.KittenName
}

func (I *InitialChecks) SetKittenName(k string) {
	I.KittenName = k
}

func (I *InitialChecks) GetDomain() string {
	return I.Domain
}

func (I *InitialChecks) SetDomain(d string) {
	I.Domain = d
}

func (I *InitialChecks) GetPid() int {
	return I.Pid
}

func (I *InitialChecks) SetPid(p int) {
	I.Pid = p
}

func (I *InitialChecks) GetPName() string {
	return I.PName
}

func (I *InitialChecks) SetPName(p string) {
	I.PName = p
}

func (I *InitialChecks) GetPath() string {
	return I.Path
}

func (I *InitialChecks) SetPath(p string) {
	I.Path = p
}
