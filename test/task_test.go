package test

import (
	"KittyStager/internal/task"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	taskTag            string
	taskPayload        []byte
	taskPayload2       []byte
	taskKey            []byte
	taskWrongKey       []byte
	taskWrongKeyLenght []byte
)

func init() {
	taskTag = "kittenSleep"
	taskPayload = []byte("this is a test payload")
	taskPayload2 = []byte("this is a test payload2")
	taskKey = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJE")
	taskWrongKey = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJ2")
	taskWrongKeyLenght = []byte("xzfbmR6MskR8J6Zr58RrhMc325kejLJ2asd")

}

func TestMarshallUnmashallTask(t *testing.T) {
	t.Parallel()
	// given
	ta := task.NewTask(taskTag, taskPayload)
	// when
	b, err := ta.MarshallTask()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = ta.UnmarshallTask(b)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, ta.GetTag(), taskTag)
	assert.Equal(t, ta.GetPayload(), taskPayload)
}

func TestEncryptDecryptTask(t *testing.T) {
	t.Parallel()
	// given
	ta := task.NewTask(taskTag, taskPayload)
	// when
	e, err := ta.EncryptTask(taskKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = ta.DecryptTask(e, taskKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, ta.GetPayload(), taskPayload)
}

func TestGetSetTask(t *testing.T) {
	t.Parallel()
	// given
	ta := task.NewTask(taskTag, taskPayload)
	assert.Equal(t, ta.GetPayload(), taskPayload)
	assert.Equal(t, ta.GetTag(), taskTag)
	ta.SetPayload(taskPayload2)
	assert.Equal(t, ta.GetPayload(), taskPayload2)
}

func TestEncryptDecryptTaskWrongKey(t *testing.T) {
	t.Parallel()
	// given
	ta := task.NewTask(taskTag, taskPayload)
	// when
	e, err := ta.EncryptTask(taskKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = ta.DecryptTask(e, taskWrongKey)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func TestEncryptTaskWrongKeyLenght(t *testing.T) {
	t.Parallel()
	// given
	ta := task.NewTask(taskTag, taskPayload)
	// when
	_, err := ta.EncryptTask(taskWrongKeyLenght)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func TestDecryptTaskWrongKeyLenght(t *testing.T) {
	t.Parallel()
	// given
	ta := task.NewTask(taskTag, taskPayload)
	// when
	e, err := ta.EncryptTask(taskKey)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = ta.DecryptTask(e, taskWrongKeyLenght)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}
