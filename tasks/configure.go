package tasks

import (
	"app/command"
	"app/uproject"
	"app/wsl"
	"github.com/akamensky/argparse"
	"os"
	"path/filepath"
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
	Context.ProjectFiles(*ctx.Opts...)
}

func (c *UE4Context) ProjectFiles(args ...string) error {
	builder := filepath.Join(c.uproject.EngineRoot, "Build", "BatchFiles", "GenerateProjectFiles.bat")
	cmdargs := make([]string, 0)
	_, err := os.Stat(builder)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		cmdargs = append(cmdargs,
			uproject.UNREAL_VERSION_SELECTOR,
			"/projectfiles")
	} else {
		cmdargs = append(cmdargs,
			`C:\Windows\System32\cmd.exe`, "/c",
			wsl.WinPath(builder))
	}
	cmdargs = append(cmdargs, wsl.WinPath(c.uproject.UProjectPath), "-Game", "-Engine", "-makefile", "-VSCode", "-2022")
	return c.run(cmdargs)
}
