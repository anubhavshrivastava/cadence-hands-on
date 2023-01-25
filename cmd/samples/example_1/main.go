package main

import (
	"flag"
)

func main() {
	var command string
	flag.StringVar(&command, "cmd", "start_workflow", "Command is start_workflow, start_worker, send_signal")
	flag.Parse()
}
