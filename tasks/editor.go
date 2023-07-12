package tasks

import (
	"app/wsl"
	"path/filepath"
	"app/command"
	"github.com/akamensky/argparse"
)

type (
	editorTarget struct {
	}
)

func InitEditor(c *command.Context) {
	command := c.Parser.NewCommand("editor", "Launch editor")
	c.Add(command, &editorTarget{})
}

func (t *editorTarget) Execute(ctx *command.Context, cmd *argparse.Command) {
	task := New()
	task.Editor(*ctx.Opts...)
}

func (c *UE4Context) Editor(args... string) error {
	editorBin := filepath.Join(c.uproject.EngineRoot, "Binaries", "Win64", "UE4Editor.exe")

	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		wsl.WinPath(editorBin),
		wsl.WinPath(c.uproject.UProjectPath),
		"-skipcompile",
		"-fullcrashdump",
		"-NOVERIFYGC")

	return c.start(cmdargs)
}
