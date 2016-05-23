package main

import (
	"flag"
	"fmt"

	"github.com/shidenggui/gotq/example/task"
)

var (
	A = flag.Int64("a", 3, "a args")
	B = flag.Int64("b", 4, "b args")
	M = flag.String("m", "request", "dispatch mode: request(block) or delay(unblock)")
)

func main() {
	flag.Parse()
	args := &task.AddArgs{*A, *B}
	if *M == "request" {
		res, err := task.AddSender.Request(args, 3)
		if err != nil {
			fmt.Println("delay error: " + err.Error())
		}
		fmt.Println(res)
		return
	}
	err := task.AddSender.Delay(args)
	if err != nil {
		fmt.Println("delay error: " + err.Error())
	}
}
