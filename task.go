package gotq

import (
	"github.com/google/uuid"
)

type Task struct {
	F     string
	Id    string
	Async bool
	Args  interface{}
}

func (t *Task) Init() {
	id, _ := uuid.NewUUID()
	t.Id = id.String()
}
