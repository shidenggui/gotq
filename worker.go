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

	}
}

// func (w *Worker) Consume()
