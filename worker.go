package gotq

import (
	"encoding/json"
	"fmt"

	log "github.com/inconshreveable/log15"
	"github.com/shidenggui/gotq/brokers"
	_ "github.com/shidenggui/gotq/log"
)

type Worker struct {
	Tasks     map[string]*TaskSender
	QueueName string
	Broker    brokers.Broker
}

var TaskChan = make(chan []byte)

func (w *Worker) Start() {
	for {
		taskStruct := new(Task)
		taskByte := <-TaskChan
		_ = json.Unmarshal(taskByte, &taskStruct)
		argsMap := taskStruct.Args
		log.Info(fmt.Sprintf("[WORKER] Receive task: %+v", taskStruct))

		sender := w.Tasks[taskStruct.F]
		res := sender.F(argsMap.(map[string]interface{}))
		log.Info(fmt.Sprintf("[WORKER] Finish task: %s, result: %v", taskStruct.Id, res))

		if taskStruct.Async {
			log.Info(fmt.Sprintf("[WORKER] Finish task: %s, result: %v, async task no need reply", taskStruct.Id, res))
			continue
		}

		log.Info(fmt.Sprintf("[WORKER] Finish task: %s, result: %v, start reply", taskStruct.Id, res))

		replyJson, err := json.Marshal(res)
		if err != nil {
			log.Info(fmt.Sprintf("[WORKER] Encode task %s result:  to json err: %s", taskStruct.Id, res, err))
			continue
		}
		err = w.Broker.Delay(replyJson, taskStruct.Id)
		if err != nil {
			log.Info(fmt.Sprintf("[WORKER] Reply task %s err: %s", taskStruct.Id, err))
			continue
		}
		log.Info(fmt.Sprintf("[WORKER] Reply task %s success", taskStruct.Id))

		w.Broker.Expire(taskStruct.Id, 1800)
	}
}

// func (w *Worker) Consume()
