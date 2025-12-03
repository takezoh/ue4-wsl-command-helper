package tasks

import (
	"app/command"
	"app/wsl"
	"github.com/akamensky/argparse"
	"path/filepath"
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
	Context.Editor(*ctx.Opts...)
}

func (c *UE4Context) Editor(args ...string) error {
	editorBin := filepath.Join(c.uproject.EngineRoot, "Binaries", "Win64", "UE4Editor.exe")
	if c.uproject.IsUE5 {
		editorBin = filepath.Join(c.uproject.EngineRoot, "Binaries", "Win64", "UnrealEditor.exe")
	}

	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		wsl.WinPath(editorBin),
		wsl.WinPath(c.uproject.UProjectPath),
		"-skipcompile",
		"-fullcrashdump",
		"-NOVERIFYGC")

	cmdargs = append(cmdargs, args...)

	return c.start(cmdargs)
}
