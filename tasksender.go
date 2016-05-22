package gotq

import (
	"encoding/json"

	"github.com/shidenggui/gotq/brokers"
)

type TaskSender struct {
	Name      string
	Broker    brokers.Broker
	F         func(map[string]interface{}) interface{}
	QueueName string
	Args      interface{}
	Result    interface{}
}

func (t *TaskSender) Delay(args interface{}) []byte {
	task := new(Task)
	task.Init()
	task.Args = args
	task.F = t.Name

	taskJson, _ := json.Marshal(task)
	t.Broker.Delay(taskJson, t.QueueName)
	return taskJson
}
