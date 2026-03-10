package command

import (
	"context"
	"io"
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

func newExecCmd(command []string) (*exec.Cmd, error) {
	println(">>>")
	println("RUN: " + strings.Join(command, " "))
	println("<<<")
	cmd := exec.Command(command[0], command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			return
		}
	}()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			return
		}
	}()

	return cmd, nil
}

func (c *Context) Run(command []string) error {
	cmd, err := newExecCmd(command)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (c *Context) Start(command []string) error {
	cmd, err := newExecCmd(command)
	if err != nil {
		return err
	}
	return cmd.Start()
}
