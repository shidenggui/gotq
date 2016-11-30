package gotq

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/shidenggui/gotq/brokers"
	"github.com/shidenggui/gotq/config"
	_ "github.com/shidenggui/gotq/log"
)

// App init broker, register task function, send task and wait result
type App struct {
	Tasks  map[string]*TaskSender
	Cfg    *config.Config
	Broker brokers.Broker
}

// Entry point for init broker
func New(cfg *config.Config) *App {
	app := new(App)
	app.Cfg = cfg
	app.Broker = brokers.NewRedisBroker(cfg.Broker)
	app.Tasks = make(map[string]*TaskSender)
	return app
}


// Register task function and return task sender
func (a *App) Register(f func(map[string]interface{}) map[string]interface{}) *TaskSender {
	taskSender := &TaskSender{
		Name:      GetFuncName(f),
		F:         f,
		Broker:    a.Broker,
		QueueName: "gotq",
	}
	a.Tasks[taskSender.Name] = taskSender
	return taskSender
}

// Start worker to consume task
func (a *App) WorkerStart(num int64) {
	for i := int64(0); i < num; i++ {
		worker := &Worker{
			Broker:    a.Broker,
			QueueName: "gotq",
			Tasks:     a.Tasks,
		}
		go worker.Start()
	}
	log.Info(fmt.Sprintf("[MainProcess] Success start workers, num: %v", num))

	brokerCfg := a.Cfg.Broker
	log.Info(fmt.Sprintf("[MainProcess] Start receive task from %s@%s:%d/%d", brokerCfg.Password, brokerCfg.Host, brokerCfg.Port, brokerCfg.DB))
	for {
		taskByte, err := a.Broker.Receive("gotq")
		if err != nil {
			log.Error(fmt.Sprintf("[MainProcess] Cant get task from redis: %#v ", err))
			time.Sleep(time.Millisecond * 100)
			continue
		}

		log.Info("[MainProcess] Block for receive task...")
		TaskChan <- taskByte
	}
}
