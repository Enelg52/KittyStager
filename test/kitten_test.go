package test

import (
	"KittyStager/internal/kitten"
	"KittyStager/internal/recon"
	"KittyStager/internal/task"
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

var (
	name      string
	sleep     int
	lastSeen  time.Time
	kittenKey string
	//Opaque     *opaque.User
)

func kittenBeforeAll() {
	name = "test"
	sleep = 5
	lastSeen = time.Now()
	kittenKey = "asdfasd"
	//Opaque     *opaque.User
}

func TestKittenSetGet(t *testing.T) {
	t.Parallel()
	// given
	kit := kitten.NewKitten(name, sleep, lastSeen, kittenKey)
	rec := recon.NewRecon(Hostname, Username, Domain, Ip, PName, Path, Pid)
	tas := task.NewTask(Tag, Payload)
	tSlice := make([]*task.Task, 0)
	taskSlice := append(tSlice, tas)
	// when
	kit.SetRecon(rec)
	kit.SetTask(tas)
	recon2 := kit.GetRecon()
	task2 := kit.GetTasks()
	// then
	assert.Equal(t, rec, recon2)
	assert.Equal(t, taskSlice[0], task2[0])
}
