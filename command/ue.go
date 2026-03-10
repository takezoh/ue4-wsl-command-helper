package command

import (
	"app/uproject"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UE struct {
	ctx      *Context
	UProject *uproject.UProject
}

func NewUE(ctx *Context, uprj *uproject.UProject) *UE {
	return &UE{ctx: ctx, UProject: uprj}
}

func (u *UE) runBuild(command string, target string, platform string, configuration string, args ...string) error {
	csproj := filepath.Join(u.UProject.EngineRoot, "Source", "Programs", target, filepath.Base(target)+".csproj")
	_, err := os.Stat(csproj)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs, `C:\Windows\System32\cmd.exe`, "/c")
	if os.IsNotExist(err) {
		command := strings.Title(command)
		cmdargs = append(cmdargs,
			filepath.Join(u.UProject.EngineRoot, "Build", "BatchFiles", command+".bat"),
			target,
			platform,
			configuration,
			u.UProject.UProjectPath)
		cmdargs = append(cmdargs, args...)
		cmdargs = append(cmdargs,
			"-unattended",
			"-fullcrashdump")
	} else {
		command := strings.ToLower(command)
		if command == "rebuild" {
			command = "build"
		}
		cmdargs = append(cmdargs,
			filepath.Join(u.UProject.EngineRoot, "Build", "BatchFiles", "MSBuild.bat"),
			"/t:"+command,
			csproj,
			"/p:GenerateFullPaths=true",
			"/p:DebugType=portable",
			"/p:Configuration="+configuration,
			"/p:Platform=AnyCPU",
			"/verbosity:minimal")
	}
	return u.ctx.Run(cmdargs)
}

func (u *UE) Build(command string, target string, platform string, configuration string, args ...string) error {
	if err := u.runBuild(command, "UnrealBuildTool", "Win64", "Development"); err != nil {
		return err
	}
	if err := u.runBuild(command, "UnrealLightmass", "Win64", "Development"); err != nil {
		return err
	}
	if err := u.runBuild(command, "ShaderCompileWorker", "Win64", "Development"); err != nil {
		return err
	}
	return u.runBuild(command, target, platform, configuration, args...)
}

func has(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func (u *UE) Package(target string, platform string, configuration string, isServer bool, args ...string) error {
	archiveName := platform + "_" + configuration + "_" + time.Now().Format("20060102_150405.00000")
	archiveDir := filepath.Join(u.UProject.ProjectRoot, "Saved", "Packages", archiveName)
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("LOGNAME")
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	cmdargs := make([]string, 0)
	unrealExe := "-ue4exe=" + u.UProject.CmdExe
	if u.UProject.IsUE5 {
		unrealExe = "-unrealexe=" + u.UProject.CmdExe
	}
	cmdargs = append(cmdargs,
		`C:\Windows\System32\cmd.exe`, "/c",
		filepath.Join(u.UProject.EngineRoot, "Build", "BatchFiles", "RunUAT.bat"),
		"-ScriptsForProject="+u.UProject.UProjectPath,
		"BuildCookRun",
		"-unattended",
		"-nocompileeditor",
		"-nop4",
		"-project="+u.UProject.UProjectPath,
		"-target="+target,
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
		"-compressed",
		"-archive", "-archivedirectory="+archiveDir,
		"-mapsonly",
		"-CrashReporter")
	if isServer {
		cmdargs = append(cmdargs, "-server")
	}
	cmdargs = append(cmdargs, args...)
	cmdargs = append(cmdargs,
		fmt.Sprintf(`-addcmdline=-statnamedevents -StatCmds='unit,fps' -SessionId=%v -SessionOwner='%v' -SessionName='%v' -messaging`,
			uuid.NewString(), username, u.UProject.Name))
	if configuration == "Shipping" {
		cmdargs = append(cmdargs, "-nodebuginfo")
	} else {
		cmdargs = append(cmdargs, "-debuginfo")
		if has("-cook", cmdargs) {
			cmdargs = append(cmdargs, "-interactivecooking")
		}
	}
	return u.ctx.Run(cmdargs)
}

func (u *UE) ProjectFiles(args ...string) error {
	builder := filepath.Join(u.UProject.EngineRoot, "Build", "BatchFiles", "GenerateProjectFiles.bat")
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
			builder)
	}
	cmdargs = append(cmdargs, u.UProject.UProjectPath, "-VSCode")
	return u.ctx.Run(cmdargs)
}

func (u *UE) Editor(args ...string) error {
	editorBin := filepath.Join(u.UProject.EngineRoot, "Binaries", "Win64", "UE4Editor.exe")
	if u.UProject.IsUE5 {
		editorBin = filepath.Join(u.UProject.EngineRoot, "Binaries", "Win64", "UnrealEditor.exe")
	}

	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		editorBin,
		u.UProject.UProjectPath,
		"-skipcompile",
		"-fullcrashdump",
		"-NOVERIFYGC")
	cmdargs = append(cmdargs, args...)
	return u.ctx.Start(cmdargs)
}

func (u *UE) Command(run string, args ...string) error {
	cmd := filepath.Join(u.UProject.EngineRoot, "Binaries", "Win64", "UE4Editor-Cmd.exe")
	if _, err := os.Stat(cmd); err != nil {
		cmd = filepath.Join(u.UProject.EngineRoot, "Binaries", "Win64", "UnrealEditor-Cmd.exe")
	}

	cmdargs := make([]string, 0)
	cmdargs = append(cmdargs,
		cmd,
		u.UProject.UProjectPath,
		"-run="+run)
	cmdargs = append(cmdargs, args...)
	return u.ctx.Run(cmdargs)
}
