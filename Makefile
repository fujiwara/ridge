.PHONY: test

test:
	go test -v ./...

install:
	go install github.com/fujiwara/ridge/cmd/ridge
