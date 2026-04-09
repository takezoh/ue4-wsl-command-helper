---
description: |
  Package (cook and archive) UE4/UE5 projects.
  Use when the user asks to package, cook, or archive an Unreal Engine project.

  Required parameters: -t <target> -p <platform> -c <configuration>
  Optional: --server, -A (additional UAT arguments)

  If the user omits parameters, find targets via `ls Source/*.Target.cs` and infer defaults from context.
argument-hint: "package [options]"
context: fork
---

# UE Package Skill

Execute `ue.exe package` based on `$ARGUMENTS`. Do NOT ask for confirmation or prompt for parameters.

1. Parse `$ARGUMENTS` to determine parameters (CLI flags or natural language)
2. If required parameters cannot be determined from `$ARGUMENTS`, return the list of missing parameters and their possible values to the main context. Do NOT execute the command

## How to run

The binary is at `../bin/ue.exe` relative to this SKILL.md. Resolve and execute:

```bash
"$SKILL_DIR/../bin/ue.exe" package [options]
```

## Execution strategy

Packaging can take tens of minutes. Use `run_in_background`:

```bash
set -o pipefail
ue.exe package -t MyGame -p Win64 -c Shipping 2>&1 | tee /tmp/ue-package-$$.log
```

## Output markers

### ue.exe markers

```
========== COMMAND STARTED ==========
RUN: <command line>
=====================================
... output ...
========== COMMAND COMPLETED ==========
```

Failure:
```
========== COMMAND FAILED ==========
ERR: <error message>
====================================
```

### Package stages

`package` runs: `BUILD -> COOK -> STAGE -> PACKAGE -> ARCHIVE`

```
********** BUILD COMMAND COMPLETED **********
********** COOK COMMAND STARTED **********
```

Use both marker types to pinpoint where a failure occurred.

## Reporting results to main context

This skill runs in a forked context. You MUST return a clear result message to the parent:

- **Success**: Report the command executed and that it completed successfully
- **Failure**: Report the command executed, the failed stage (BUILD/COOK/STAGE/PACKAGE/ARCHIVE), and the relevant error lines from the output
- **Insufficient parameters**: Do NOT guess or prompt interactively. Return the list of missing parameters and their possible values to the main context so the user can clarify
