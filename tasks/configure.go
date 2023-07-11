package tasks

import (
	"app/command"
	"os"
	"path"
	"github.com/akamensky/argparse"
)

type (
	configureTarget struct {
	}
)

func InitConfigure(c *command.Context) {
	command := c.Parser.NewCommand("configure", "Make build configuration")
	c.Add(command, &configureTarget{})
}

func (t *configureTarget) Execute(ctx *command.Context, cmd *argparse.Command) {
	task := New()
	task.ProjectFiles(*ctx.Opts...)
}

func (c *UE4Context) ProjectFiles(args... string) error {
	builder := path.Join(c.uproject.EngineRoot, "Build", "BatchFiles", "GenerateProjectFiles.bat")
	cmdargs := make([]string, 0)
	_, err := os.Stat(builder)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		cmdargs = append(cmdargs,
			c.upath(UNREAL_VERSION_SELECTOR),
			"/projectfiles")
	} else {
		cmdargs = append(cmdargs,
			c.upath("C:/Windows/System32/cmd.exe"),
			"/c",
			c.wpath(builder))
	}
	cmdargs = append(cmdargs, c.wpath(c.uproject.UProjectPath), "-Game", "-Engine", "-makefile", "-VSCode")
	return c.run(cmdargs)
}
