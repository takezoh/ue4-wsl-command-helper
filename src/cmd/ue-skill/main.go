package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/takezoh/ue-cli-skill/command"
	"github.com/takezoh/ue-cli-skill/tasks"
	"github.com/takezoh/ue-cli-skill/uproject"
)

func main() {
	sub := "ue"
	if len(os.Args) > 1 {
		sub = os.Args[1]
	}

	f, err := os.CreateTemp("", "ue-"+sub+"-*.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(os.Stdout, "LOG_FILE: %s\n", f.Name())

	stdout := io.MultiWriter(os.Stdout, f)
	stderr := io.MultiWriter(os.Stderr, f)
	ctx := command.New(context.Background(), stdout, stderr)

	uprj, err := uproject.GetUProject()
	if err != nil {
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}
	if err = os.Chdir(uprj.RootPath); err != nil {
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(stdout, "Found: "+uprj.UProjectPath)
	fmt.Fprintln(stdout, "Set workspace: "+uprj.RootPath)

	ue := command.NewUE(ctx, uprj)
	p := command.NewParser(ue)

	tasks.InitConfigure(p, ue)
	tasks.InitBuilds(p, ue)
	tasks.InitEditor(p, ue)
	tasks.InitCommand(p, ue)

	if err := p.Parse(os.Args); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}
