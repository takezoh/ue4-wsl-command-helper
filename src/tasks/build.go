package tasks

import (
	"github.com/takezoh/ue-cli-skill/command"

	"github.com/akamensky/argparse"
)

type buildTask struct {
	opts          *[]string
	platform      *string
	configuration *string
	target        *string
	isServer      *bool
}

func InitBuilds(p *command.Parser, ue *command.UE) {
	add := func(name string, desc string) {
		cmd := p.ArgParser.NewCommand(name, desc)
		t := &buildTask{opts: p.Opts}
		t.target = cmd.Selector("t", "target", ue.UProject.Targets, &argparse.Options{Required: true, Help: "Build target"})
		t.platform = cmd.Selector("p", "platform", ue.UProject.Platforms, &argparse.Options{Required: true, Help: "Target platform"})
		t.configuration = cmd.Selector("c", "configuration", ue.UProject.Configurations, &argparse.Options{Required: true, Help: "Target configuration"})
		if name == "package" {
			t.isServer = cmd.Flag("", "server", &argparse.Options{Required: false, Help: "Build to be server role"})
		}
		p.Add(cmd, t)
	}
	add("build", "Build")
	add("clean", "Clean")
	add("rebuild", "Rebuild")
	add("package", "Make package")
}

func (t *buildTask) Do(ue *command.UE, cmd *argparse.Command) error {
	switch cmd.GetName() {
	case "build", "clean", "rebuild":
		return ue.Build(cmd.GetName(), *t.target, *t.platform, *t.configuration, *t.opts...)
	case "package":
		return ue.Package(*t.target, *t.platform, *t.configuration, *t.isServer, *t.opts...)
	}
	return nil
}
