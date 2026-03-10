all:
	GOOS=windows GOARCH=amd64 go build -o bin/ue.exe .
