package tasks

import (
	"github.com/takezoh/ue-cli-skill/command"

	"github.com/akamensky/argparse"
)

type commandTask struct {
	opts *[]string
	run  *string
}

func InitCommand(p *command.Parser, ue *command.UE) {
	cmd := p.ArgParser.NewCommand("command", "Run commandlet")
	t := &commandTask{opts: p.Opts}
	t.run = cmd.String("r", "run", &argparse.Options{Required: true, Help: "Run target"})
	p.Add(cmd, t)
}

func (t *commandTask) Do(ue *command.UE, cmd *argparse.Command) error {
	return ue.Command(*t.run, *t.opts...)
}
