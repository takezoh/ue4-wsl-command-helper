---
context: fork
---

# UE Build Skill

Run Unreal Engine 4/5 build and development commands.

## Usage

`ue.exe` is a Windows binary. It can be called directly from WSL via interop.

The binary is located at `bin/ue.exe` relative to this SKILL.md.
Resolve the path from the location of this file.

Example: if this SKILL.md is at `~/.claude/skills/ue/SKILL.md`:
```bash
~/.claude/skills/ue/bin/ue.exe <command> [options]
```

If project-local (`.claude/skills/ue/SKILL.md`):
```bash
.claude/skills/ue/bin/ue.exe <command> [options]
```

## Commands

```bash
# Build
ue.exe build -t <target> -p <platform> -c <configuration>
ue.exe clean -t <target> -p <platform> -c <configuration>
ue.exe rebuild -t <target> -p <platform> -c <configuration>

# Package
ue.exe package -t <target> -p <platform> -c <configuration> [--server]

# Generate project files
ue.exe configure

# Launch editor
ue.exe editor

# Run commandlet
ue.exe command -r <commandlet>
```

## Arguments

### Finding targets

```bash
ls Source/*.Target.cs
```

The target name is the filename without the `.Target.cs` extension.

### platform

`Win64`, `Android`, `Linux`, `PS4`, `PS5`

### configuration

`Development`, `Test`, `Shipping`

### Extra options

`-A` / `--opt` passes additional arguments directly to UE tools.

## Environment variables

- `TARGET_DIR`: Explicitly specify the project directory (defaults to auto-detecting `.uproject` from cwd)

## Notes

- Builds can take tens of minutes. Set `timeout: 600000` on Bash tool calls
- For builds expected to exceed 10 minutes, prefer background execution:
  ```bash
  ue.exe build -t MyGame -p Win64 -c Development > /tmp/ue-build.log 2>&1 &
  ```
