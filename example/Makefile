export AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text)

handler: *.go
	GOOS=linux GOARCH=amd64 go build -o handler main.go

deploy: handler function.json
	lambroll deploy

clean:
	rm -f handler

logs:
	lambroll logs --follow
