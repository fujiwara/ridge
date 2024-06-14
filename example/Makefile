export AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text)

bootstrap: *.go
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

deploy: bootstrap function.json
	lambroll deploy

clean:
	rm -f bootstrap

logs:
	lambroll logs --follow
