package uproject

import (
	"app/wsl"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	CONFIGURATION_LIST = []string{"Development", "Test", "Shipping"}
	PLATFORM_LIST      = []string{"Win64", "Android", "Linux", "PS4", "PS5"}
)

const (
	UNREAL_VERSION_SELECTOR    = "C:/Program Files (x86)/Epic Games/Launcher/Engine/Binaries/Win64/UnrealVersionSelector.exe"
	UNREAL_ENGINE_INSTALL_ROOT = "C:/Program Files/Epic Games"
)

type (
	UProject struct {
		IsUE5        bool

		Name         string
		ProjectRoot  string
		RootPath     string
		EngineRoot   string
		UProjectPath string
		CmdExe    string

		Targets        []string
		Configurations []string
		Platforms      []string

		Modules    []*Module
		HasModules bool

		FileVersion       int
		EngineAssociation string
	}
	Module struct {
		Name                   string
		Type                   string
		LoadingPhase           string
		AdditionalDependencies []string
	}
)

func GetUProject() (*UProject, error) {
	path, err := uprojectPath()
	if err != nil {
		return nil, err
	}

	uprj, err := uprojectObj(path)
	if err != nil {
		return nil, err
	}
	return uprj, nil
}

func uprojectPath() (string, error) {
	currentDir := os.Getenv("TARGET_DIR")
	if currentDir == "" {
		var err error
		currentDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	for currentDir != string(os.PathSeparator) {
		matches, _ := filepath.Glob(filepath.Join(currentDir, "*.uproject"))
		if len(matches) == 0 {
			matches, _ = filepath.Glob(filepath.Join(currentDir, "*", "*.uproject"))
		}
		for _, v := range matches {
			if !strings.HasSuffix(v, "EngineTest.uproject") {
				return v, nil
			}
		}
		currentDir = filepath.Dir(currentDir)
	}

	return "", errors.New("Not Found *.uproject")
}

func uprojectObj(uprojectPath string) (*UProject, error) {
	prj := new(UProject)

	f, err := os.Open(uprojectPath)
	if err != nil {
		return nil, err
	}

	b, _ := io.ReadAll(f)
	if bytes.Equal(b[:3], []byte{0xef, 0xbb, 0xbf}) {
		b = b[3:]
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return nil, err
	}
	prj.FileVersion = int(obj["FileVersion"].(float64))
	prj.EngineAssociation = obj["EngineAssociation"].(string)
	prj.Modules = parseModules(obj["Modules"].([]interface{}))
	prj.HasModules = len(prj.Modules) > 0

	prj.UProjectPath = uprojectPath
	prj.Name = filepath.Base(uprojectPath[:len(uprojectPath)-len(".uproject")])
	prj.ProjectRoot = filepath.Dir(uprojectPath)
	prj.RootPath = filepath.Dir(prj.ProjectRoot)
	{
		editorpathMatches, _ := filepath.Glob(filepath.Join(prj.RootPath, ".ue4-version"))
		if len(editorpathMatches) == 0 {
			editorpathMatches, _ = filepath.Glob(filepath.Join(prj.RootPath, "*", ".ue4-version"))
		}
		if len(editorpathMatches) > 0 {
			prj.RootPath = filepath.Dir(editorpathMatches[0])
			if b, err := os.ReadFile(editorpathMatches[0]); err != nil {
				prj.EngineRoot = string(b)
			}
		} else {
			enginePathMatches, _ := filepath.Glob(filepath.Join(prj.RootPath, "Engine", "Build", "Build.version"))
			if len(enginePathMatches) == 0 {
				enginePathMatches, _ = filepath.Glob(filepath.Join(prj.RootPath, "*", "Engine", "Build", "Build.version"))
			}
			prj.EngineRoot = filepath.Dir(filepath.Dir(enginePathMatches[0]))
		}
		if _, err := os.Stat(prj.EngineRoot); os.IsNotExist(err) {
			prj.RootPath = prj.ProjectRoot
			prj.EngineRoot = filepath.Join(wsl.UnixPath(UNREAL_ENGINE_INSTALL_ROOT), "UE_"+prj.EngineAssociation, "Engine")
		}
	}
	prj.IsUE5 = false
	prj.CmdExe = filepath.Join(prj.EngineRoot, "Binaries", "Win64", "UE4Editor-Cmd.exe")
	if _, err := os.Stat(prj.CmdExe); err != nil {
		prj.IsUE5 = true
		prj.CmdExe = filepath.Join(prj.EngineRoot, "Binaries", "Win64", "UnrealEditor-Cmd.exe")
	}

	prj.Targets = targets(prj.ProjectRoot)
	prj.Configurations = CONFIGURATION_LIST
	prj.Platforms = PLATFORM_LIST
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

func targets(projectRoot string) []string {
	matches, _ := filepath.Glob(filepath.Join(projectRoot, "Source", "*.Target.cs"))
	targets := []string{}
	for _, v := range matches {
		t := filepath.Base(v[:len(v)-len(".Target.cs")])
		targets = append(targets, t)
	}
	return targets
}
