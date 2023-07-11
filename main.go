package main

import (
	"app/tasks"
	"app/command"
	"os"
)

func main() {
	cmd := command.New()

	tasks.InitConfigure(cmd)
	tasks.InitBuilds(cmd)
	tasks.InitEditor(cmd)

	err := cmd.Parse(os.Args)
	if err != nil {
		os.Exit(1)
	}
}