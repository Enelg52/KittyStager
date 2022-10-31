package httpUtil

import (
	"KittyStager/cmd/util"
	"time"
)

type Kitten struct {
	Name       string
	Payload    []byte
	Sleep      int
	Id         int
	LastSeen   time.Time
	InitChecks util.InitialChecks
	Alive      bool
}

func (K *Kitten) GetTarget() string {
	return K.Name
}

func (K *Kitten) GetPayload() []byte {
	return K.Payload
}

func (K *Kitten) SetPayload(sc []byte) {
	K.Payload = sc
}

func (K *Kitten) GetId() int {
	return K.Id
}

func (K *Kitten) SetId(id int) {
	K.Id = id
}

func (K *Kitten) GetLastSeen() time.Time {
	return K.LastSeen
}

func (K *Kitten) SetLastSeen(t time.Time) {
	K.LastSeen = t
}

func (K *Kitten) GetInitChecks() util.InitialChecks {
	return K.InitChecks
}

func (K *Kitten) SetInitChecks(c util.InitialChecks) {
	K.InitChecks = c
}

func (K *Kitten) GetSleep() int {
	return K.Sleep
}

func (K *Kitten) SetSleep(t int) {
	K.Sleep = t
}

func (K *Kitten) GetAlive() bool {
	return K.Alive
}

func (K *Kitten) SetAlive(b bool) {
	K.Alive = b
}
