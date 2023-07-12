package tasks

import (
	"app/uproject"
	"app/wsl"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

var Context *UE4Context

type (
	UE4Context struct {
		ctx      context.Context
		uproject *uproject.UProject
	}
)

func init() {
	Context = new(UE4Context)
	Context.ctx = context.Background()

	uprj, err := uproject.GetUProject()
	if err != nil {
		panic(err)
	}
	Context.uproject = uprj

	if err = os.Chdir(uprj.RootPath); err != nil {
		panic(err)
	}

	println("Found: " + uprj.UProjectPath)
	println("Set workspace: " + uprj.RootPath)
}

func newExecCmd(command []string) (*exec.Cmd, error) {
	println(">>>")
	println("RUN: " + strings.Join(command, " "))
	println("<<<")
	cmd := exec.Command(wsl.UnixPath(command[0]), command[1:]...)

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
