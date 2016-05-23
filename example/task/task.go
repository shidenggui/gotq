package task

import (
	"fmt"

	ms "github.com/mitchellh/mapstructure"
	"github.com/shidenggui/gotq"
	"github.com/shidenggui/gotq/config"
)

// gotq config
var cfg = config.Config{
	Broker: &config.BrokerCfg{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		DB:       0,
	},
}
var App = gotq.New(&cfg)

// taskSender
var AddSender = App.Register(Add)

type AddArgs struct {
	X int64
	Y int64
}

type AddResult struct {
	Sum int64
}

func Add(argsInter map[string]interface{}) map[string]interface{} {
	args := new(AddArgs)
	ms.Decode(argsInter, args)

	fmt.Println("args ", args)

	res := make(map[string]interface{})
	res["sum"] = args.X + args.Y
	return res
}
