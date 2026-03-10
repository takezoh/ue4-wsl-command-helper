all:
	cd src && GOOS=windows GOARCH=amd64 go build -o ../skills/ue/bin/ue.exe .
