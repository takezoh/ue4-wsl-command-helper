package command

import (
	"fmt"

	"github.com/akamensky/argparse"
)

type (
	Parser struct {
		ArgParser *argparse.Parser
		Opts      *[]string

		ue      *UE
		targets []taskWrapper
	}

	taskWrapper struct {
		cmd  *argparse.Command
		task Task
	}

	Task interface {
		Do(ue *UE, cmd *argparse.Command)
	}
)

func (p *Parser) Add(cmd *argparse.Command, task Task) {
	p.targets = append(p.targets, taskWrapper{cmd: cmd, task: task})
}

func (p *Parser) Parse(args []string) error {
	err := p.ArgParser.Parse(args)
	if err != nil {
		fmt.Print(p.ArgParser.Usage(err))
		return err
	}
	for _, v := range p.targets {
		if v.cmd.Happened() {
			v.task.Do(p.ue, v.cmd)
		}
	}
	return nil
}

func NewParser(ue *UE) *Parser {
	p := &Parser{ue: ue}
	p.ArgParser = argparse.NewParser("ue", "Execute ue4 commands")
	p.Opts = p.ArgParser.StringList("A", "opt", &argparse.Options{Required: false, Help: "Options"})
	return p
}
