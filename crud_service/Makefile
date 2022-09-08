.PHONY: run
run:
	go run cmd/bot/main.go

build:
	go build -o bin/bot cmd/bot/main.go

build-linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/bot-linux64 cmd/bot/main.go

build-win64:
	GOOS=windows GOARCH=amd64 go build -o bin/bot-linux64 cmd/bot/main.go

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

MIGRATIONS_DIR=./migrations
.PHONY: migration
migration:
	goose -dir=${MIGRATIONS_DIR} create $(NAME) sql

.PHONY: .test
.test:
	$(info Running tests...)
	go test ./...

.PHONY: cover
cover:
	go test -v $$(go list ./... | grep -v -E './pkg/(api)') -covermode=count -coverprofile=./c.out
	go tool cover -html=./c.out
