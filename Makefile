.PHONY: test test-rie

test:
	go test -v ./...

# requires Docker
test-rie:
	go test -v -tags rie -run TestRIE .

install:
	go install github.com/fujiwara/ridge/cmd/ridge
