package test

import (
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"KittyStager/internal/task/recon"
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

var (
	kittenName   string
	kittenSleep  int
	lastSeen     time.Time
	kittenAlive  bool
	kittenKey    string
	kittenRecon  *recon.Recon
	kittenResult *task.Task

	kittenTask *task.Task
)

func init() {
	kittenName = "test"
	kittenSleep = 5
	lastSeen = time.Now()
	kittenAlive = true
	kittenKey = "asdfasd"
	kittenTask = task.NewTask("test", []byte("tasks"))
	kittenRecon = recon.NewRecon("host", "user", "domain", "ip", "pName", "path", 1000)
	kittenResult = task.NewTask("test", []byte("kittenResult"))
}

func TestKittenSetGet(t *testing.T) {
	t.Parallel()
	// given
	kit := kitten.NewKitten(kittenName, kittenSleep, lastSeen, kittenKey)

	// when
	kit.SetTask(kittenTask)
	kit.SetRecon(kittenRecon)
	kit.SetSleep(kittenSleep)
	kit.SetLastSeen(lastSeen)
	kit.SetAlive(kittenAlive)
	kit.SetResult(kittenResult)
	// then
	assert.Equal(t, kit.GetRecon(), kittenRecon)
	tasks := kit.GetTasks()
	assert.Equal(t, tasks[0], kittenTask)
	assert.Equal(t, kit.GetAlive(), kittenAlive)
	assert.Equal(t, kit.GetLastSeen(), lastSeen)
	assert.Equal(t, kit.GetSleep(), kittenSleep)
	assert.Equal(t, kit.GetKey(), kittenKey)
	assert.Equal(t, kit.GetName(), kittenName)
	assert.Equal(t, kit.GetResult(), kittenResult)
}
