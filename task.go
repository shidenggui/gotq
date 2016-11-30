package gotq

import (
	"github.com/google/uuid"
)

// Broker send task struct
type Task struct {
	// Task invoke function name
	F        string
	// Task uuid, contain create timestamp
	Id       string
	// Whether to return result
	Async    bool
	// Task function args
	Args     interface{}
	// Indicate block seconds when Async is True
	WaitTime int64
}

// Generate task function name, args, and whether need return result
func (t *Task) Init(f string, args interface{}, async bool) {
	id, _ := uuid.NewUUID()
	t.Id = id.String()

	t.F = f
	t.Args = args
	t.Async = async
}
