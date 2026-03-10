package command

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

type Context struct {
	ctx context.Context
}

func New(ctx context.Context) *Context {
	return &Context{ctx: ctx}
}

func newExecCmd(ctx context.Context, command []string) *exec.Cmd {
	println("========== COMMAND STARTED ==========")
	println("RUN: " + strings.Join(command, " "))
	println("=====================================")
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (c *Context) Run(command []string) error {
	err := newExecCmd(c.ctx, command).Run()
	if err != nil {
		println("========== COMMAND FAILED ==========")
		println("ERR: " + err.Error())
		println("====================================")
	} else {
		println("========== COMMAND COMPLETED ==========")
	}
	return err
}

func (c *Context) Start(command []string) error {
	return newExecCmd(c.ctx, command).Start()
}
