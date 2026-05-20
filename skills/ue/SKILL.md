---
description: |
  Run UE4/UE5 build, package, configure, editor, and commandlet commands.
  Use when the user asks to build, cook, package, clean, or launch UE projects.

  Subcommands and required parameters:
    build/clean/rebuild: -t <target> -p <platform> -c <configuration>
    package: -t <target> -p <platform> -c <configuration> [--server]
    configure: (no parameters)
    editor: (no parameters)
    command: -r <commandlet>

  If the user omits parameters, find targets via `ls Source/*.Target.cs` and infer defaults from context.
argument-hint: "<subcommand> [options]"
---

# UE Build Skill

Execute `ue.exe` based on `$ARGUMENTS`. Do NOT ask for confirmation or prompt for parameters.

1. Parse `$ARGUMENTS` to determine the subcommand and parameters (CLI flags or natural language)
2. If required parameters cannot be determined from `$ARGUMENTS`, return the list of missing parameters and their possible values to the main context. Do NOT execute the command

## How to run

The binary is at `bin/ue.exe` relative to this SKILL.md. Resolve and execute:

```bash
/path/to/skill/bin/ue.exe <subcommand> [options]
```

## Execution strategy

Run `ue-skill.exe` (located next to `ue.exe`). It tees all output to a temp log file
and prints the path on stdout before any other output.

```bash
/path/to/skill/bin/ue-skill.exe build -t MyGame -p Win64 -c Development
```

Builds can take tens of minutes — use `run_in_background: true` for `build` and `package`.

## Detecting success / failure

1. **Exit code is authoritative.** Non-zero = failure (`ue-skill.exe` propagates the child process exit code unchanged).
2. **Find the log.** The first line of stdout is `LOG_FILE: <path>`. Open that file for the full output.
3. **Locate the failed stage** (package only): grep the log for `********** <STAGE> COMMAND COMPLETED **********` markers to determine which of `BUILD / COOK / STAGE / PACKAGE / ARCHIVE` failed.

## Output markers (in the log file)

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

### UE package stages

`package` runs: `BUILD -> COOK -> STAGE -> PACKAGE -> ARCHIVE`

```
********** BUILD COMMAND COMPLETED **********
********** COOK COMMAND STARTED **********
```

## Reporting results to main context

This skill runs in a forked context. You MUST return a clear result message to the parent:

- **Success**: Report the command executed and that it completed successfully
- **Failure**: Report the command executed, the failed stage (for package: BUILD/COOK/STAGE/PACKAGE/ARCHIVE), and the relevant error lines from the output
- **Insufficient parameters**: Do NOT guess or prompt interactively. Return the list of missing parameters and their possible values to the main context so the user can clarify
