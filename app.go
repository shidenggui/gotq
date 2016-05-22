package gotq

import (
	"github.com/shidenggui/gotq/brokers"
	"github.com/shidenggui/gotq/config"
)

type App struct {
	Tasks  map[string]*TaskSender
	Cfg    *config.Config
	Broker brokers.Broker
}

func New(cfg *config.Config) *App {
	app := new(App)
	app.Cfg = cfg
	app.Broker = brokers.NewRedisBroker(cfg.Broker)
	app.Tasks = make(map[string]*TaskSender)
	return app
}

func (a *App) Register(f func(map[string]interface{}) map[string]interface{}, args interface{}, result interface{}) *TaskSender {
	taskSender := &TaskSender{
		Name:      GetFuncName(f),
		F:         f,
		Broker:    a.Broker,
		QueueName: "gotq",
		Args:      args,
		Result:    result,
	}
	a.Tasks[taskSender.Name] = taskSender
	return taskSender
}

func (a *App) WorkerStart(num int64) {
	for i := int64(0); i < num; i++ {
		worker := &Worker{
			Broker:    a.Broker,
			QueueName: "gotq",
			Tasks:     a.Tasks,
		}
		go worker.Start()
	}

	for {
		taskByte, err := a.Broker.Receive("gotq")
		if err != nil {
			panic(err)
		}
		TaskChan <- taskByte
	}
}
