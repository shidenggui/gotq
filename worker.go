package gotq

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
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
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error(fmt.Sprintf("[WORKER] recover from panic: %v", r))
					debug.PrintStack()
				}
			}()
			taskStruct := new(Task)
			taskByte := <-TaskChan
			_ = json.Unmarshal(taskByte, &taskStruct)
			argsMap := taskStruct.Args
			log.Info(fmt.Sprintf("[WORKER] Receive task: %+v", taskStruct))

			// drop sync outdate task
			if !taskStruct.Async {
				uuidStruct, err := uuid.Parse(taskStruct.Id)
				if err != nil {
					log.Error("[WORKER] task id error, cant parse uuid")
					return
				}
				uuidTime := uuidStruct.Time()
				sendTime, _ := uuidTime.UnixTime()

				nowTime := time.Now()
				interval := nowTime.Unix() - sendTime
				if interval > taskStruct.WaitTime {
					log.Info(fmt.Sprintf("[WORKER] discard task outdate, task %#v, send time: %#v, wait time: %#v, now time: %#v", taskStruct, sendTime, taskStruct.WaitTime, time.Now()))
					return
				}

			}

			sender, ok := w.Tasks[taskStruct.F]
			if !ok {
				log.Error(fmt.Sprintf("[WORKER] Task %s not register, register task list: %#v", taskStruct.F, w.Tasks))
				return
			}
			res := sender.F(argsMap.(map[string]interface{}))

			if taskStruct.Async {
				log.Info(fmt.Sprintf("[WORKER] Finish task: %s, result: %v, async task no need reply", taskStruct.Id, res))
				return
			}

			log.Info(fmt.Sprintf("[WORKER] Finish task: %s, reply result: %v", taskStruct.Id, res))

			replyJson, err := json.Marshal(res)
			if err != nil {
				log.Info(fmt.Sprintf("[WORKER] Encode task %s result:  to json err: %s", taskStruct.Id, res, err))
				return
			}
			err = w.Broker.Delay(replyJson, taskStruct.Id)
			if err != nil {
				log.Info(fmt.Sprintf("[WORKER] Reply task %s err: %s", taskStruct.Id, err))
				return
			}

			w.Broker.Expire(taskStruct.Id, 1800)
		}()
	}
}
