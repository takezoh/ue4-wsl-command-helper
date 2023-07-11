package tasks

import (
	"strings"
	"os/exec"
)

func (c *UE4Context) wslpath(path string, opt string) string {
	if !c.isWSL {
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

func (c *UE4Context) upath(upath string) string {
	return c.wslpath(upath, "-au")
}

func (c *UE4Context) wpath(upath string) string {
	return c.wslpath(upath, "-am")
}
