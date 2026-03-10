# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

WSL2 上で Unreal Engine 4/5 のビルド・開発作業を実行するコマンドラインツール。`.uproject` を自動検出してプロジェクトコンテキストを初期化し、Windows バッチスクリプトを WSL から呼び出す。

## Commands

```bash
# ビルド実行（直接）
go run ./main.go build -t <target> -p <platform> -c <configuration>
go run ./main.go clean -t <target> -p <platform> -c <configuration>
go run ./main.go rebuild -t <target> -p <platform> -c <configuration>
go run ./main.go package -t <target> -p <platform> -c <configuration> [--server]
go run ./main.go configure
go run ./main.go editor
go run ./main.go command -r <commandlet>

# Mage タスク経由
mage Build <target> <platform> <configuration>
mage Clean <target> <platform> <configuration>
mage Rebuild <target> <platform> <configuration>
mage Package <platform> <configuration>
mage Configure
mage Editor
```

Platform: `Win64, Android, Linux, PS4, PS5`
Configuration: `Development, Test, Shipping`

テスト・lint コマンドは未定義。

## Architecture

### 実行フロー

```
main.go
  → command.New()          # コマンドフレームワーク初期化
  → Init*(ctx) x4          # タスク登録 (Configure/Build/Editor/Command)
  → cmd.Parse(os.Args)     # 引数解析・対応タスク実行
```

### 主要コンポーネント

- **`command/command.go`**: `Target` インターフェース（`Execute(ctx, cmd)`）と argparse ベースのディスパッチャ
- **`tasks/context.go`**: `UE4Context` グローバル変数。`init()` で `.uproject` を探索してプロジェクトルートに `cd`、以降すべてのタスクがこのコンテキストを参照
- **`uproject/uproject.go`**: `.uproject` JSON 解析。UE4/UE5 自動判別（コマンド実行ファイル名で判定）、エンジンルート解決、ターゲット抽出
- **`wsl/wsl.go`**: `wslpath` コマンドで Unix ↔ Windows パス変換（`UnixPath()` / `WinPath()`）
- **`tasks/build.go`**: UnrealBuildTool 呼び出し。C# プロジェクトは MSBuild.bat、ネイティブは Build.bat を使い分け
- **`tasks/command.go`**: Commandlet 実行（`UE4Editor-Cmd.exe` / `UnrealEditor-Cmd.exe`）

### UE4/UE5 分岐

`uproject.go` でエンジンバージョンを判別し、`UE4Context` に保持。各タスクで実行ファイル名を切り替える（`UE4Editor.exe` vs `UnrealEditor.exe` など）。

### 追加オプション

`-A` / `--opt` フラグで任意の追加引数を UE ツールにそのまま渡せる。
