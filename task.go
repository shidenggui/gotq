package gotq

import (
	"github.com/google/uuid"
)

type Task struct {
	F        string
	Id       string
	Async    bool
	Args     interface{}
	WaitTime int64
}

func (t *Task) Init(f string, args interface{}, async bool) {
	id, _ := uuid.NewUUID()
	t.Id = id.String()

	t.F = f
	t.Args = args
	t.Async = async
}
