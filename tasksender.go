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


// Just send task to broker, doesn't return result
//
// If you need task result, should use Request
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

// Just send task to broker, doesn't return result
//
// When broker is redis, it'll use rpush to send task, Delay default use lpush.
// This can cause task send later but execute quicker
//
// If you need task result should use Request
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

// Send task to broker and wait result
//
// Return result is type of map[string]interface{}, must binding by yourself
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

	// block for reply
	replyByte, err := t.Broker.Request(task.Id, blockTime)
	if err != nil {
		return nil, err
	}
	replyMap := make(map[string]interface{})
	json.Unmarshal(replyByte, &replyMap)
	return replyMap, nil
}

// Send task to broker and wait result
//
// When broker is redis, it'll use rpush to send task, Request default use lpush.
// This can cause task send later but execute quicker
//
// Return result is type of map[string]interface{}, must binding by yourself
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
