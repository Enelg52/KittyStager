package util

type Task struct {
	Tag     string `json:"Tag"`
	Payload []byte `json:"Payload"`
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
