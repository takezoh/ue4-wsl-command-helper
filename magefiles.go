//go:build mage
// +build mage
package main

import (
	"app/tasks"
)

func Hoge() {
	tasks.New()
	println("hoge")
}

func Configure() error {
	t := tasks.New()
	return t.ProjectFiles()
}

func Projectfiles() error {
	t := tasks.New()
	return t.ProjectFiles()
}

// func Build(target string, command string, configuration string, platform string, args... string) error {
func Build(target string, platform string, configuration string) error {
	t := tasks.New()
	// return t.Build(target, command, configuration, platform, args...)
	return t.Build("build", target, platform, configuration)
}
func Clean(target string, platform string, configuration string) error {
	t := tasks.New()
	return t.Build("clean", target, platform, configuration)
}
func Rebuild(target string, platform string, configuration string) error {
	t := tasks.New()
	return t.Build("rebuild", target, platform, configuration)
}
func Package(platform string, configuration string) error {
	t := tasks.New()
	return t.Package(platform, configuration)
}

func Editor() error {
	t := tasks.New()
	return t.Editor()
}
