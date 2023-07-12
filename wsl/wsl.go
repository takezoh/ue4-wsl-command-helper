package wsl

import (
	"os"
	"os/exec"
	"strings"
)

var isWSL = os.Getenv("WSL_DISTRO_NAME") != ""

func Wslpath(path string, opt string) string {
	if !isWSL {
		return path
	}
	if path == "" {
		return ""
	}

	cmd := exec.Command("wslpath", opt, path)
	b, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(b))
}

func UnixPath(upath string) string {
	return Wslpath(upath, "-au")
}

func WinPath(upath string) string {
	return Wslpath(upath, "-aw")
}
