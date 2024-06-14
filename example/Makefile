export AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text)

bootstrap: *.go
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

deploy: bootstrap function.json
	lambroll deploy --function-url function_url.json

clean:
	rm -f bootstrap

logs:
	lambroll logs --follow
