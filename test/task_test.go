package test

import (
	"KittyStager/internal/task"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	Tag     string
	Payload []byte
)

func taskBeforeAll() {
	Tag = "sleep"
	Payload = []byte("Il a test payload")

}

func TestMarshallUnmashallTask(t *testing.T) {
	t.Parallel()
	// given
	task := task.NewTask(Tag, Payload)
	// when
	b, err := task.MarshallTask()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = task.UnmarshallTask(b)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, task.GetTag(), Tag)
	assert.Equal(t, task.GetPayload(), Payload)
}

func TestEncryptDecryptTask(t *testing.T) {
	t.Parallel()
	// given
	task := task.NewTask(Tag, Payload)
	// when
	e, err := task.EncryptTask(Key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = task.DecryptTask(e, Key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, task.GetPayload(), Payload)
}
