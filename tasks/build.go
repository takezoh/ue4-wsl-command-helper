package tasks

import (
	"app/command"
	"app/wsl"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	buildTarget struct {
		platform      *string
		configuration *string
		target        *string
		isServer      *bool
	}
)

func InitBuilds(c *command.Context) {
	c.Add(newBuildCmd(c.Parser.NewCommand("build", "Build")))
	c.Add(newBuildCmd(c.Parser.NewCommand("clean", "Clean")))
	c.Add(newBuildCmd(c.Parser.NewCommand("rebuild", "Rebuild")))
	c.Add(newBuildCmd(c.Parser.NewCommand("package", "Make package")))
}

func newBuildCmd(command *argparse.Command) (*argparse.Command, *buildTarget) {
	t := new(buildTarget)
	t.target = command.Selector("t", "target", Context.uproject.Targets, &argparse.Options{Required: true, Help: "Build target"})
	t.platform = command.Selector("p", "platform", Context.uproject.Platforms, &argparse.Options{Required: true, Help: "Target platform"})
	t.configuration = command.Selector("c", "configuration", Context.uproject.Configurations, &argparse.Options{Required: true, Help: "Target configuration"})
	if command.GetName() == "package" {
		t.isServer = command.Flag("", "server", &argparse.Options{Required: false, Help: "Build to be server role"})
	}
	return command, t
}

func (t *buildTarget) Execute(ctx *command.Context, cmd *argparse.Command) {
	switch cmd.GetName() {
	case "build", "clean", "rebuild":
		Context.Build(cmd.GetName(), *t.target, *t.platform, *t.configuration, *ctx.Opts...)
	case "package":
		Context.Package(*t.target, *t.platform, *t.configuration, *t.isServer, *ctx.Opts...)
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

func (c *UE4Context) runBuild(command string, target string, platform string, configuration string, args ...string) error {
	csproj := filepath.Join(c.uproject.EngineRoot, "Source", "Programs", target, filepath.Base(target)+".csproj")
	_, err := os.Stat(csproj)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs, `C:\Windows\System32\cmd.exe`, "/c")
	if os.IsNotExist(err) {
		command := strings.Title(command)
		cmdargs = append(cmdargs,
			wsl.WinPath(filepath.Join(c.uproject.EngineRoot, "Build", "BatchFiles", command+".bat")),
			target,
			platform,
			configuration,
			wsl.WinPath(c.uproject.UProjectPath))
		cmdargs = append(cmdargs, args...)
		cmdargs = append(cmdargs,
			// "-verbose",
			"-unattended",
			"-fullcrashdump")
	} else {
		command := strings.ToLower(command)
		if command == "rebuild" {
			command = "build"
		}
		cmdargs = append(cmdargs,
			wsl.WinPath(filepath.Join(c.uproject.EngineRoot, "Build", "BatchFiles", "MSBuild.bat")),
			"/t:"+command,
			wsl.WinPath(csproj),
			"/p:GenerateFullPaths=true",
			"/p:DebugType=portable",
			"/p:Configuration="+configuration,
			"/p:Platform=AnyCPU",
			"/verbosity:minimal")
	}
	return c.run(cmdargs)
}

func (c *UE4Context) Build(command string, target string, platform string, configuration string, args ...string) error {
	// c.runBuild(command, "DotNETCommon/DotNETUtilities", platform, "Development")
	// c.runBuild(command, "UnrealHeaderTool", "Win64", "Development")
	c.runBuild(command, "UnrealBuildTool", "Win64", "Development")
	// c.runBuild(command, "AutomationTool", "Win64", "Development")
	c.runBuild(command, "UnrealLightmass", "Win64", "Development")
	c.runBuild(command, "ShaderCompileWorker", "Win64", "Development")
	return c.runBuild(command, target, platform, configuration, args...)
}

func (c *UE4Context) Package(target string, platform string, configuration string, isServer bool, args ...string) error {
	archiveName := platform + "_" + configuration + "_" + time.Now().Format("20060102_150405.00000")
	archiveDir := filepath.Join(c.uproject.ProjectRoot, "Saved", "Packages", archiveName)
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("LOGNAME")
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	cmdargs := make([]string, 0)
	unrealExe := "-ue4exe="+wsl.WinPath(c.uproject.CmdExe)
	if c.uproject.IsUE5 {
		unrealExe = "-unrealexe="+wsl.WinPath(c.uproject.CmdExe)
	}
	cmdargs = append(cmdargs,
		`C:\Windows\System32\cmd.exe`, "/c",
		wsl.WinPath(filepath.Join(c.uproject.EngineRoot, "Build", "BatchFiles", "RunUAT.bat")),
		"-ScriptsForProject="+wsl.WinPath(c.uproject.UProjectPath),
		"BuildCookRun",
		"-unattended",
		"-nocompileeditor",
		"-nop4",
		"-project="+wsl.WinPath(c.uproject.UProjectPath),
		"-target="+target,
		// '-SkipCookingEditorContent',
		"-clientconfig="+configuration,
		"-serverconfig="+configuration,
		"-targetplatform="+platform,
		unrealExe,
		"-ddc=DerivedDataBackendGraph",
		"-utf8output",
		"-fullcrashdump")
	cmdargs = append(cmdargs,
		"-cookflavor=ASTC",
		"-prereqs",
		"-build", "-compile",
		"-cook", "-stage",
		"-pak", "-package",
		// "-distribution",
		// "-nodebuginfo"
		"-compressed",
		"-archive", "-archivedirectory="+wsl.WinPath(archiveDir),
		"-mapsonly",
		"-CrashReporter")
	if isServer {
		cmdargs = append(cmdargs, "-server")
	}
	cmdargs = append(cmdargs, args...)
	cmdargs = append(cmdargs,
		fmt.Sprintf(`-addcmdline=-statnamedevents -StatCmds='unit,fps' -SessionId=%v -SessionOwner='%v' -SessionName='%v' -messaging`,
			uuid.NewString(), username, c.uproject.Name))
	if configuration == "Shipping" {
		cmdargs = append(cmdargs, "-nodebuginfo")
	} else {
		cmdargs = append(cmdargs, "-debuginfo")
		if has("-cook", cmdargs) {
			cmdargs = append(cmdargs, "-interactivecooking")
		}
	}
	return c.run(cmdargs)
}
