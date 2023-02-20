package task

import (
	"KittyStager/pkg/crypto"
	"encoding/json"
)

type Task struct {
	Tag     string `json:"Tag"`
	Payload []byte `json:"Payload"`
}

func NewTask(tag string, payload []byte) *Task {
	return &Task{Tag: tag, Payload: payload}
}

func (task *Task) EncryptTask(key []byte) ([]byte, error) {
	c := crypto.NewChaCha20()
	m, err := task.MarshallTask()
	if err != nil {
		return nil, err
	}
	e, err := c.Encrypt(m, key)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (task *Task) DecryptTask(cypherText []byte, key []byte) error {
	c := crypto.NewChaCha20()
	d, err := c.Decrypt(cypherText, key)
	if err != nil {
		return err
	}
	err = task.UnmarshallTask(d)
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) MarshallTask() ([]byte, error) {
	t, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (task *Task) UnmarshallTask(j []byte) error {
	err := json.Unmarshal(j, &task)
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) GetTag() string {
	return task.Tag
}

func (task *Task) GetPayload() []byte {
	return task.Payload
}

func (task *Task) SetPayload(p []byte) {
	task.Payload = p
}
