package tasks

import (
	"app/command"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/akamensky/argparse"
)

type (
	buildTarget struct {
		platform *string
		configuration *string
		target *string
	}
)

func InitBuilds(c *command.Context) {
	c.Add(newBuildCmd(c.Parser.NewCommand("build", "Build")))
	c.Add(newBuildCmd(c.Parser.NewCommand("clean", "Clean")))
	c.Add(newBuildCmd(c.Parser.NewCommand("rebuild", "Rebuild")))
	c.Add(newPackageCmd(c.Parser.NewCommand("package", "Make package")))
}

func newBuildCmd(command *argparse.Command) (*argparse.Command, *buildTarget) {
	t := new(buildTarget)
	t.platform = command.String("p", "platform", &argparse.Options{Required: true, Help: "Target platform"})
	t.configuration = command.String("c", "configuration", &argparse.Options{Required: true, Help: "Target configuration"})
	t.target = command.String("t", "target", &argparse.Options{Required: true, Help: "Build target"})
	return command, t
}

func newPackageCmd(command *argparse.Command) (*argparse.Command, *buildTarget) {
	t := new(buildTarget)
	t.platform = command.String("p", "platform", &argparse.Options{Required: true, Help: "Target platform"})
	t.configuration = command.String("c", "configuration", &argparse.Options{Required: true, Help: "Target configuration"})
	return command, t
}

func (t *buildTarget) Execute(ctx *command.Context, cmd *argparse.Command) {
	task := New()
	switch cmd.GetName() {
	case "build", "clean", "rebuild": task.Build(cmd.GetName(), *t.target, *t.platform, *t.configuration, *ctx.Opts...)
	case "package": task.Package(*t.platform, *t.configuration, *ctx.Opts...)
	}
}

func has(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func (c *UE4Context) runBuild(command string, target string, platform string, configuration string, args... string) error {
	csproj := path.Join(c.uproject.EngineRoot, "Source", "Programs", target, path.Base(target)+".csproj")
	_, err := os.Stat(csproj)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		c.upath("C:/Windows/System32/cmd.exe"),
		"/c")
	if os.IsNotExist(err) {
		command := strings.Title(command)
		cmdargs = append(cmdargs,
			c.wpath(path.Join(c.uproject.EngineRoot, "Build", "BatchFiles", command+".bat")),
			target,
			platform,
			configuration)
		cmdargs = append(cmdargs, args...)
		cmdargs = append(cmdargs,
			// "-verbose",
			"-fullcrashdump")
	} else {
		command := strings.ToLower(command)
		if command == "rebuild" {
			command = "build"
		}
		cmdargs = append(cmdargs,
			c.wpath(path.Join(c.uproject.EngineRoot, "Build", "BatchFiles", "MSBuild.bat")),
			"/t:"+command,
			c.wpath(csproj),
			"/p:GenerateFullPaths=true",
			"/p:DebugType=portable",
			"/p:Configuration="+configuration,
			"/p:Platform=AnyCPU",
			"/verbosity:minimal")
	}
	return c.run(cmdargs)
}

func (c *UE4Context) Build(command string, target string, platform string, configuration string, args... string) error {
	// c.runBuild(command, "DotNETCommon/DotNETUtilities", platform, "Development")
	// c.runBuild(command, "UnrealHeaderTool", "Win64", "Development")
	// c.runBuild(command, "UnrealBuildTool", "Win64", "Development")
	c.runBuild(command, "AutomationTool", "Win64", "Development")
	c.runBuild(command, "UnrealLightmass", "Win64", "Development")
	c.runBuild(command, "ShaderCompileWorker", "Win64", "Development")
	return c.runBuild(command, target, platform, configuration, args...)
}

func (c *UE4Context) Package(platform string, configuration string, args... string) error {
	archiveName := platform +"_"+ configuration +"_"+ time.Now().Format("20060102_150405.00000")
	archiveDir := path.Join(c.uproject.ProjectRoot, "Saved", "Packages", archiveName)
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("LOGNAME")
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		c.upath("C:/Windows/System32/cmd.exe"), "/c",
		c.wpath(path.Join(c.uproject.EngineRoot, "Build", "BatchFiles", "RunUAT.bat")),
		"-ScriptsForProject="+c.wpath(c.uproject.UProjectPath),
		"BuildCookRun",
		"-nocompileeditor",
		"-nop4",
		"-project="+c.wpath(c.uproject.UProjectPath),
		// '-SkipCookingEditorContent',
		"-clientconfig="+configuration,
		"-prereqs", "-targetplatform="+platform, "-utf8output",
		"-fullcrashdump")
	cmdargs = append(cmdargs,
		"-cookflavor=ASTC",
		"-build", "-compile",
		"-cook", "-stage",
		"-pak",
		"-package",
		// "-distribution",
		// "-nodebuginfo"
		"-compressed",
		"-archive",
		"-archivedirectory="+c.wpath(archiveDir),
		"-mapsonly")
		//+' '+ (opts or ''))
	cmdargs = append(cmdargs,
		"-serverconfig="+configuration,
		fmt.Sprintf(`-addcmdline=-statnamedevents -StatCmds='unit,fps' -SessionId=%v -SessionOwner='%v' -SessionName='%v' -messaging`,
			uuid.NewString(), username, c.uproject.Name))
	if configuration == "Shipping" {
		cmdargs = append(cmdargs, "-nodebuginfo")
	} else {
		cmdargs = append(cmdargs, "-debuginfo", "-CrashReporter")
		if has("-cook", cmdargs) {
			cmdargs = append(cmdargs, "-interactivecooking")
		}
	}
	return c.run(cmdargs)
}
