# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Windows binary (`ue.exe`) for running Unreal Engine 4/5 build and development tasks. Auto-detects `.uproject` to initialize project context and invokes Windows batch scripts. Can be called directly from WSL via interop.

## Build

Cross-compile as a Windows binary (eliminates the need for wslpath conversion when called via WSL interop).

```bash
make        # produces .claude/skills/ue/bin/ue.exe
```

## Commands

```bash
.claude/skills/ue/bin/ue.exe build -t <target> -p <platform> -c <configuration>
.claude/skills/ue/bin/ue.exe clean -t <target> -p <platform> -c <configuration>
.claude/skills/ue/bin/ue.exe rebuild -t <target> -p <platform> -c <configuration>
.claude/skills/ue/bin/ue.exe package -t <target> -p <platform> -c <configuration> [--server]
.claude/skills/ue/bin/ue.exe configure
.claude/skills/ue/bin/ue.exe editor
.claude/skills/ue/bin/ue.exe command -r <commandlet>
```

Platform: `Win64, Android, Linux, PS4, PS5`
Configuration: `Development, Test, Shipping`

No test or lint commands defined.

## Architecture

### Execution flow

```
main.go
  -> command.New()          # Initialize command framework
  -> Init*(ctx) x4          # Register tasks (Configure/Build/Editor/Command)
  -> cmd.Parse(os.Args)     # Parse args and dispatch to task
```

### Key components

- **`command/command.go`**: `Target` interface (`Execute(ctx, cmd)`) and argparse-based dispatcher
- **`tasks/context.go`**: `UE4Context` global. `init()` searches for `.uproject` and `cd`s to project root; all tasks reference this context
- **`uproject/uproject.go`**: `.uproject` JSON parsing. Auto-detects UE4/UE5 (by executable name), resolves engine root, extracts targets
- **`tasks/build.go`**: UnrealBuildTool invocation. Uses MSBuild.bat for C# projects, Build.bat for native
- **`tasks/command.go`**: Commandlet execution (`UE4Editor-Cmd.exe` / `UnrealEditor-Cmd.exe`)

### UE4/UE5 branching

`uproject.go` detects the engine version and stores it in `UE4Context`. Each task switches executable names accordingly (`UE4Editor.exe` vs `UnrealEditor.exe`, etc.).

### Extra options

`-A` / `--opt` flag passes arbitrary additional arguments directly to UE tools.

## Skill installation

To use as a Claude Code Skill in other UE projects:

```bash
mkdir -p .claude/skills && ln -s /path/to/ue4-wsl-command-helper/.claude/skills/ue .claude/skills/ue
```
