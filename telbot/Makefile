run:
	go run cmd/bot/main.go

build:
	go build -o bin/bot cmd/bot/main.go

build-linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/bot-linux64 -ldflags "-s -w" cmd/bot/main.go

build-win64:
	GOOS=windows GOARCH=amd64 go build -o bin/bot-linux64 cmd/bot/main.go