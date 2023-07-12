package command

import (
	"fmt"
	"github.com/akamensky/argparse"
)

type (
	Context struct {
		Parser *argparse.Parser
		Opts   *[]string

		targets []targetWrapper
	}

	targetWrapper struct {
		cmd    *argparse.Command
		target Target
	}

	Target interface {
		Execute(ctx *Context, cmd *argparse.Command)
	}
)

func (c *Context) Add(cmd *argparse.Command, target Target) {
	c.targets = append(c.targets, targetWrapper{cmd: cmd, target: target})
}

func (c *Context) Parse(args []string) error {
	err := c.Parser.Parse(args)
	if err != nil {
		fmt.Print(c.Parser.Usage(err))
		return err
	}
	for _, v := range c.targets {
		if v.cmd.Happened() {
			v.target.Execute(c, v.cmd)
		}
	}
	return nil
}

func New() *Context {
	ctx := new(Context)
	ctx.Parser = argparse.NewParser("ue", "Execute ue4 commands")
	ctx.Opts = ctx.Parser.StringList("A", "opt", &argparse.Options{Required: false, Help: "Options"})
	return ctx
}
