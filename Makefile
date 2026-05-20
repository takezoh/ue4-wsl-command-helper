BIN := skills/ue/bin

all: $(BIN)/ue.exe $(BIN)/ue-skill.exe

$(BIN)/ue.exe:
	cd src && GOOS=windows GOARCH=amd64 go build -o ../$(BIN)/ue.exe ./cmd/ue

$(BIN)/ue-skill.exe:
	cd src && GOOS=windows GOARCH=amd64 go build -o ../$(BIN)/ue-skill.exe ./cmd/ue-skill

.PHONY: all
