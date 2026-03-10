package tasks

import (
	"github.com/takezoh/ue-cli-skill/command"

	"github.com/akamensky/argparse"
)

type configureTask struct {
	opts *[]string
}

func InitConfigure(p *command.Parser, ue *command.UE) {
	cmd := p.ArgParser.NewCommand("configure", "Make build configuration")
	p.Add(cmd, &configureTask{opts: p.Opts})
}

func (t *configureTask) Do(ue *command.UE, cmd *argparse.Command) {
	ue.ProjectFiles(*t.opts...)
}
