all:
	GOOS=windows GOARCH=amd64 go build -o .claude/skills/ue/bin/ue.exe .
