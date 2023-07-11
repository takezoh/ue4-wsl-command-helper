package tasks

import (
	"app/uproject"
	"os"
	"os/exec"
	"io"
	"context"
	"strings"
)

type (
	UE4Context struct {
		ctx context.Context
		uproject *uproject.UProject
	}
)

func New() *UE4Context {
	ctx := new(UE4Context)
	ctx.ctx = context.Background()

	uprj, err := uproject.GetUProject()
	if err != nil {
		panic(err)
	}
	ctx.uproject = uprj

	if err = os.Chdir(uprj.RootPath); err != nil {
		panic(err)
	}

	println("Found: "+uprj.UProjectPath)
	println("Set workspace: "+uprj.RootPath)
	return ctx
}

func newExecCmd(command []string) (*exec.Cmd, error) {
	println(">>> RUN: ")
	println(strings.Join(command, " "))
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

func (c *UE4Context) run(command []string) error {
	cmd, err := newExecCmd(command)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (c *UE4Context) start(command []string) error {
	cmd, err := newExecCmd(command)
	if err != nil {
		return err
	}
	return cmd.Start()
}
