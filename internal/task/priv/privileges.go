package priv

import (
	"encoding/json"
	"github.com/fourcorelabs/wintoken"
)

type Privileges struct {
	Priv []wintoken.Privilege `json:"process"`
	Integrity string `json:"integrity"`
}


func NewPrivileges(priv []wintoken.Privilege, integrity string) *Privileges {
	return &Privileges{Priv: priv, Integrity: integrity}
}


func (privileges *Privileges) MarshallPrivileges() ([]byte, error) {
	t, err := json.Marshal(privileges)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (privileges *Privileges) UnmarshallPrivileges(j []byte) error {
	err := json.Unmarshal(j, &privileges)
	if err != nil {
		return err
	}
	return nil
}

