package tasks

import (
	"app/command"

	"github.com/akamensky/argparse"
)

type editorTask struct {
	opts *[]string
}

func InitEditor(p *command.Parser, ue *command.UE) {
	cmd := p.ArgParser.NewCommand("editor", "Launch editor")
	p.Add(cmd, &editorTask{opts: p.Opts})
}

func (t *editorTask) Do(ue *command.UE, cmd *argparse.Command) {
	ue.Editor(*t.opts...)
}
