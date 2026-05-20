package main

import (
	"context"
	"os"

	"github.com/takezoh/ue-cli-skill/command"
	"github.com/takezoh/ue-cli-skill/tasks"
	"github.com/takezoh/ue-cli-skill/uproject"
)

func main() {
	ctx := command.New(context.Background(), os.Stdout, os.Stderr)

	uprj, err := uproject.GetUProject()
	if err != nil {
		panic(err)
	}
	if err = os.Chdir(uprj.RootPath); err != nil {
		panic(err)
	}
	println("Found: " + uprj.UProjectPath)
	println("Set workspace: " + uprj.RootPath)

	ue := command.NewUE(ctx, uprj)
	p := command.NewParser(ue)

	tasks.InitConfigure(p, ue)
	tasks.InitBuilds(p, ue)
	tasks.InitEditor(p, ue)
	tasks.InitCommand(p, ue)

	if err := p.Parse(os.Args); err != nil {
		os.Exit(1)
	}
}
