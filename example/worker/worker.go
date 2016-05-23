package main

import (
	"flag"

	"github.com/shidenggui/gotq/example/task"
)

var (
	WorkerNum = flag.Int64("n", 500, "worker number")
)

func main() {
	flag.Parse()

	task.App.WorkerStart(10)
}
