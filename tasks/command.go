package tasks

import (
	"app/command"
	"app/wsl"
	"github.com/akamensky/argparse"
	"path/filepath"
)

type (
	commandTarget struct {
		run *string
	}
)

func InitCommand(c *command.Context) {
	command := c.Parser.NewCommand("command", "Run commandlet")
	t := commandTarget{}
	t.run = command.String("r", "run", &argparse.Options{Required: true, Help: "Run target"})
	c.Add(command, &t)
}

func (t *commandTarget) Execute(ctx *command.Context, cmd *argparse.Command) {
	Context.Command(*t.run, *ctx.Opts...)
}

func (c *UE4Context) Command(run string, args ...string) error {
	cmd := filepath.Join(c.uproject.EngineRoot, "Binaries", "Win64", "UE4Editor-Cmd.exe")

	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		wsl.WinPath(cmd),
		wsl.WinPath(c.uproject.UProjectPath),
		"-run=" + run)
	cmdargs = append(cmdargs, args...)
	return c.run(cmdargs)
}
