package tasks

import (
	"os"
	"strings"
	"context"
	"errors"
	"path"
	"path/filepath"
	"encoding/json"
)

const (
	UNREAL_VERSION_SELECTOR = "C:/Program Files (x86)/Epic Games/Launcher/Engine/Binaries/Win64/UnrealVersionSelector.exe"
	UNREAL_ENGINE_INSTALL_ROOT = "C:/Program Files/Epic Games"
)

type (
	UE4Context struct {
		ctx context.Context
		uproject *UProject
		isWSL bool
	}
	UProject struct {
		FileVersion int
		EngineAssociation string
		Modules []*Module
		HasModules bool
		UProjectPath string
		Name string
		ProjectRoot string
		RootPath string
		EngineRoot string
	}
	Module struct {
		Name string
		Type string
		LoadingPhase string
		AdditionalDependencies []string
	}
)

func New() *UE4Context {
	ctx := new(UE4Context)
	ctx.ctx = context.Background()
	ctx.isWSL = os.Getenv("WSL_DISTRO_NAME") != ""

	uproject, _ := findUproject()
	uprj, err := newUproject(ctx, uproject)
	if os.IsNotExist(err) {
		panic("*.uproject is not found.")
	}
	if err != nil {
		panic(err)
	}
	ctx.uproject = uprj

	if err = os.Chdir(ctx.uproject.RootPath); err != nil {
		panic(err)
	}

	println("Found: "+uproject)
	println("Set workspace: "+ctx.uproject.RootPath)
	return ctx
}

func newUproject(c *UE4Context, uprojectPath string) (*UProject, error) {
	prj := new(UProject)

	f, err := os.Open(uprojectPath)
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	err = json.NewDecoder(f).Decode(&obj)
	if err != nil {
		return nil, err
	}
	prj.FileVersion = int(obj["FileVersion"].(float64))
	prj.EngineAssociation = obj["EngineAssociation"].(string)
	prj.Modules = parseModules(obj["Modules"].([]interface{}))
	prj.HasModules = len(prj.Modules) > 0

	prj.UProjectPath = uprojectPath
	prj.Name = path.Base(uprojectPath[:len(uprojectPath)-len(".uproject")])
	prj.ProjectRoot = path.Dir(uprojectPath)
	prj.RootPath = path.Dir(prj.ProjectRoot)

	{
		editorpathMatches, _ := filepath.Glob(path.Join(prj.RootPath, ".ue4-version"))
		if len(editorpathMatches) == 0 {
			editorpathMatches, _ = filepath.Glob(path.Join(prj.RootPath, "*", ".ue4-version"))
		}
		if len(editorpathMatches) > 0 {
			prj.RootPath = path.Dir(editorpathMatches[0])
			if b, err := os.ReadFile(editorpathMatches[0]); err != nil {
				prj.EngineRoot = string(b)
			}
		} else {
			enginePathMatches, _ := filepath.Glob(path.Join(prj.RootPath, "Engine", "Build", "Build.version"))
			if len(enginePathMatches) == 0 {
				enginePathMatches, _ = filepath.Glob(path.Join(prj.RootPath, "*", "Engine", "Build", "Build.version"))
			}
			prj.EngineRoot = path.Dir(path.Dir(enginePathMatches[0]))
		}
		if _, err := os.Stat(prj.EngineRoot); os.IsNotExist(err) {
			prj.RootPath = prj.ProjectRoot
			prj.EngineRoot = path.Join(c.upath(UNREAL_ENGINE_INSTALL_ROOT), "UE_"+prj.EngineAssociation, "Engine")
		}
	}
	return prj, nil
}

func parseModules(src []interface{}) []*Module {
	var mods []*Module
	for _, v_ := range src {
		v := v_.(map[string]interface{})
		m := new(Module)
		m.Name = v["Name"].(string)
		m.Type = v["Type"].(string)
		m.LoadingPhase = v["LoadingPhase"].(string)
		mods = append(mods, m)
	}
	return mods
}

func findUproject() (string, error) {
	currentDir := os.Getenv("TARGET_DIR")
	if currentDir == "" {
		var err error
		currentDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	for currentDir!= "/" {
		matches, _ := filepath.Glob(path.Join(currentDir, "*.uproject"))
		if len(matches) == 0 {
			matches, _ = filepath.Glob(path.Join(currentDir, "*", "*.uproject"))
		}
		for _, v := range matches {
			if !strings.HasSuffix(v, "EngineTest.uproject") {
				return v, nil
			}
		}
		currentDir = path.Dir(currentDir)
	}

	return "", errors.New("Not Found *.uproject")
}
