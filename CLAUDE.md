# CLAUDE.md

## Build

```bash
make        # produces bin/ue.exe
```

No test or lint commands defined.

## Architecture

```
main.go
  -> command.New()
  -> Init*(ctx) x4          # Configure/Build/Editor/Command
  -> cmd.Parse(os.Args)
```

- **`command/command.go`**: `Target` interface and argparse dispatcher
- **`tasks/context.go`**: `UE4Context` global; `init()` finds `.uproject` and sets project root
- **`uproject/uproject.go`**: `.uproject` parsing, UE4/UE5 detection, engine root resolution
- **`tasks/build.go`**: UnrealBuildTool / MSBuild.bat invocation
- **`tasks/command.go`**: Commandlet execution

UE4/UE5 branching is determined in `uproject.go` by executable name. Each task switches binaries accordingly (`UE4Editor.exe` vs `UnrealEditor.exe`).
