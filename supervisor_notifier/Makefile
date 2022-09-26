run:
	go run cmd/spv_notif/main.go

build:
	go build -o bin/spv_notif cmd/spv_notif/main.go

build-linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/spv_notif_linux64 -ldflags "-s -w" cmd/spv_notif/main.go

build-win64:
	GOOS=windows GOARCH=amd64 go build -o bin/spv_notif_win64 cmd/spv_notif/main.go