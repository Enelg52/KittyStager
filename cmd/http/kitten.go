package http

import (
	"KittyStager/cmd/util"
	"github.com/frekui/opaque"
	"time"
)

// Kitten is a struct that contains all the information about a kittens
type Kitten struct {
	Name       string
	Payload    []byte
	Sleep      int
	Id         int
	LastSeen   time.Time
	InitChecks util.InitialChecks
	Alive      bool
	Key        string
	Opaque     *opaque.User
}

func NewKitten() *Kitten {
	return &Kitten{}
}

func (K *Kitten) GetName() string {
	return K.Name
}

func (K *Kitten) SetName(name string) {
	K.Name = name
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

func (K *Kitten) GetKey() string {
	return K.Key
}

func (K *Kitten) SetKey(key string) {
	K.Key = key
}

func (K *Kitten) GetOpaque() *opaque.User {
	return K.Opaque
}

func (K *Kitten) SetOpaque(u *opaque.User) {
	K.Opaque = u
}
