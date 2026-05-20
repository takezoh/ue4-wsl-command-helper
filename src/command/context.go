package command

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Context struct {
	ctx    context.Context
	Stdout io.Writer
	Stderr io.Writer
}

func New(ctx context.Context, stdout, stderr io.Writer) *Context {
	return &Context{ctx: ctx, Stdout: stdout, Stderr: stderr}
}

func (c *Context) newExecCmd(command []string) *exec.Cmd {
	fmt.Fprintln(c.Stderr, "========== COMMAND STARTED ==========")
	fmt.Fprintln(c.Stderr, "RUN: "+strings.Join(command, " "))
	fmt.Fprintln(c.Stderr, "=====================================")
	cmd := exec.CommandContext(c.ctx, command[0], command[1:]...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd
}

func (c *Context) Run(command []string) error {
	err := c.newExecCmd(command).Run()
	if err != nil {
		fmt.Fprintln(c.Stderr, "========== COMMAND FAILED ==========")
		fmt.Fprintln(c.Stderr, "ERR: "+err.Error())
		fmt.Fprintln(c.Stderr, "====================================")
	} else {
		fmt.Fprintln(c.Stderr, "========== COMMAND COMPLETED ==========")
	}
	return err
}

func (c *Context) Start(command []string) error {
	return c.newExecCmd(command).Start()
}
