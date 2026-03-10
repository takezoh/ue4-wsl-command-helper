# ue-cli-skill

Windows binary (`ue.exe`) for running Unreal Engine 4/5 build and development tasks. Auto-detects `.uproject` to initialize project context and invokes Windows batch scripts. Can be called directly from WSL via interop.

## Build

```bash
make        # produces .claude/skills/ue/bin/ue.exe
```

## Commands

```bash
ue.exe build -t <target> -p <platform> -c <configuration>
ue.exe clean -t <target> -p <platform> -c <configuration>
ue.exe rebuild -t <target> -p <platform> -c <configuration>
ue.exe package -t <target> -p <platform> -c <configuration> [--server]
ue.exe configure
ue.exe editor
ue.exe command -r <commandlet>
```

- Platform: `Win64`, `Android`, `Linux`, `PS4`, `PS5`
- Configuration: `Development`, `Test`, `Shipping`
- Extra options: `-A` / `--opt` passes additional arguments directly to UE tools

### Finding targets

```bash
ls Source/*.Target.cs
```

The target name is the filename without the `.Target.cs` extension.

## Claude Code Skill

Install as a Claude Code Skill in any UE project:

```bash
mkdir -p .claude/skills && ln -s /path/to/ue4-wsl-command-helper/.claude/skills/ue .claude/skills/ue
```

## Architecture

```
main.go
  -> command.New()          # Initialize command framework
  -> Init*(ctx) x4          # Register tasks (Configure/Build/Editor/Command)
  -> cmd.Parse(os.Args)     # Parse args and dispatch to task
```

- **`command/command.go`**: `Target` interface (`Execute(ctx, cmd)`) and argparse-based dispatcher
- **`tasks/context.go`**: `UE4Context` global. `init()` searches for `.uproject` and `cd`s to project root
- **`uproject/uproject.go`**: `.uproject` JSON parsing. Auto-detects UE4/UE5, resolves engine root, extracts targets
- **`tasks/build.go`**: UnrealBuildTool invocation. Uses MSBuild.bat for C# projects, Build.bat for native
- **`tasks/command.go`**: Commandlet execution (`UE4Editor-Cmd.exe` / `UnrealEditor-Cmd.exe`)

UE4/UE5 is auto-detected in `uproject.go` and each task switches executable names accordingly.
