package util

import "encoding/json"

// Task is the struct that contains the tasks
type Task struct {
	Tag     string `json:"Tag"`
	Payload []byte `json:"Payload"`
}

func NewTask() *Task {
	return &Task{}
}

// TaskUnmarshalJSON unmarshal the json
func TaskUnmarshalJSON(j []byte) (*Task, error) {
	iniCheck := NewTask()
	err := json.Unmarshal(j, &iniCheck)
	if err != nil {
		return &Task{}, err
	}
	return iniCheck, nil
}

func (T *Task) GetTag() string {
	return T.Tag
}

func (T *Task) SetTag(t string) {
	T.Tag = t
}

func (T *Task) GetPayload() []byte {
	return T.Payload
}

func (T *Task) SetPayload(p []byte) {
	T.Payload = p
}
