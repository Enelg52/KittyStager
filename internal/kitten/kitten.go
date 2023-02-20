package kitten

import (
	"KittyStager/internal/recon"
	"KittyStager/internal/task"
	"time"
)

type Kitten struct {
	Name     string
	Sleep    int
	LastSeen time.Time
	Alive    bool
	Key      string
	Tasks    []*task.Task
	Recon    *recon.Recon
}

func NewKitten(name string, sleep int, lastSeen time.Time, key string) *Kitten {
	return &Kitten{
		Name:     name,
		Sleep:    sleep,
		LastSeen: lastSeen,
		Alive:    true,
		Key:      key,
		Tasks:    nil,
		Recon:    nil,
	}
}

func (kitten *Kitten) SetTask(task *task.Task) {
	kitten.Tasks = append(kitten.Tasks, task)
}

func (kitten *Kitten) GetTasks() []*task.Task {
	return kitten.Tasks
}

func (kitten *Kitten) SetRecon(recon *recon.Recon) {
	kitten.Recon = recon
}

func (kitten *Kitten) GetRecon() *recon.Recon {
	return kitten.Recon
}

func (kitten *Kitten) GetAlive() bool {
	return kitten.Alive
}

func (kitten *Kitten) SetAlive(state bool) {
	kitten.Alive = state
}

func (kitten *Kitten) GetLastSeen() time.Time {
	return kitten.LastSeen
}

func (kitten *Kitten) SetLastSeen(t time.Time) {
	kitten.LastSeen = t
}

func (kitten *Kitten) GetSleep() int {
	return kitten.Sleep
}
