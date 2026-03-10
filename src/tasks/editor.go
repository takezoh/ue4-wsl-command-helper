package tasks

import (
	"github.com/takezoh/ue-cli-skill/command"

	"github.com/akamensky/argparse"
)

type editorTask struct {
	opts *[]string
}

func InitEditor(p *command.Parser, ue *command.UE) {
	cmd := p.ArgParser.NewCommand("editor", "Launch editor")
	p.Add(cmd, &editorTask{opts: p.Opts})
}

func (t *editorTask) Do(ue *command.UE, cmd *argparse.Command) error {
	return ue.Editor(*t.opts...)
}
