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
		sender := w.Tasks[taskStruct.F]
		//realArgs.X = 1111
		log.Info(fmt.Sprintf("%#v", argsMap))
		sender.F(argsMap.(map[string]interface{}))
	}
}

// func (w *Worker) Consume()
