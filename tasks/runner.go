package tasks

import (
	"strings"
	"io"
	"os"
	"os/exec"
)

func newCmd(command []string) (*exec.Cmd, error) {
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
	cmd, err := newCmd(command)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (c *UE4Context) start(command []string) error {
	cmd, err := newCmd(command)
	if err != nil {
		return err
	}
	return cmd.Start()
}
