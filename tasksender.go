package gotq

import (
	"encoding/json"

	"github.com/shidenggui/gotq/brokers"
)

type TaskSender struct {
	Name      string
	Broker    brokers.Broker
	F         func(map[string]interface{}) map[string]interface{}
	QueueName string
	Args      interface{}
}

func (t *TaskSender) Delay(args interface{}) error {
	const async = true
	task := new(Task)
	task.Init(t.Name, args, async)

	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = t.Broker.Delay(taskJson, t.QueueName)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskSender) QuickDelay(args interface{}) error {
	const async = true
	task := new(Task)
	task.Init(t.Name, args, async)

	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = t.Broker.QuickDelay(taskJson, t.QueueName)
	if err != nil {
		return err
	}
	return nil
}

// request, use redis's lpush to send tasks
func (t *TaskSender) Request(args interface{}, blockTime int64) (map[string]interface{}, error) {
	const async = false
	task := new(Task)
	task.Init(t.Name, args, async)
	task.WaitTime = blockTime

	taskJson, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	err = t.Broker.Delay(taskJson, t.QueueName)
	if err != nil {
		return nil, err
	}

	//block for reply
	replyByte, err := t.Broker.Request(task.Id, blockTime)
	if err != nil {
		return nil, err
	}
	replyMap := make(map[string]interface{})
	json.Unmarshal(replyByte, &replyMap)
	return replyMap, nil
}

// quick request, use redis's rpush to send tasks
func (t *TaskSender) QuickRequest(args interface{}, blockTime int64) (map[string]interface{}, error) {
	const async = false
	task := new(Task)
	task.Init(t.Name, args, async)
	task.WaitTime = blockTime

	taskJson, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	err = t.Broker.QuickDelay(taskJson, t.QueueName)
	if err != nil {
		return nil, err
	}

	//block for reply
	replyByte, err := t.Broker.Request(task.Id, blockTime)
	if err != nil {
		return nil, err
	}
	replyMap := make(map[string]interface{})
	json.Unmarshal(replyByte, &replyMap)
	return replyMap, nil
}
