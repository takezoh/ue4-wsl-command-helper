---
context: fork
---

# UE Build Skill

Run Unreal Engine 4/5 build and development commands.

## Usage

`ue.exe` is a Windows binary. It can be called directly from WSL via interop.

The binary is located at `bin/ue.exe` relative to this SKILL.md.
Resolve the path from the location of this file.

Example: if this SKILL.md is at `.claude/skills/ue/SKILL.md` (symlinked):
```bash
SKILL_DIR="$(dirname "$(readlink -f "$0")")"
"$SKILL_DIR/bin/ue.exe" <command> [options]
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

## Execution strategy

Builds can take tens of minutes. To avoid blocking the main context:

1. Generate a unique log file path: `/tmp/ue-build-<pid>.log`
2. Launch the build in the background:
   ```bash
   ue.exe build -t MyGame -p Win64 -c Development > /tmp/ue-build-$$.log 2>&1 &
   echo "PID=$! LOG=/tmp/ue-build-$$.log"
   ```
3. Return immediately to the main context with:
   - The exact command that was launched
   - The PID and log file path
   - How to check status: `kill -0 <pid> 2>/dev/null && echo running || echo done` and `tail -50 <log>`

## Reporting results

This skill runs in a forked context. You MUST return a clear result message to the parent context:

- **For long-running commands** (build, clean, rebuild, package): Launch in background and return the PID + log path immediately
- **For quick commands** (configure, editor, command): Run directly and report success/failure with exit code and last 50 lines of output on error
