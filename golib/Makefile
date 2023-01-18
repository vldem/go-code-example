.PHONY: .test
.test:
	$(info Running tests...)
	go test ./...

.PHONY: cover
cover:
	go test -v $$(go list ./... | grep -v -E './pkg/(api)') -covermode=count -coverprofile=./c.out
	go tool cover -html=./c.out