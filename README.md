# ridge

AWS Lambda HTTP Proxy integration event bridge to Go's net/http.

## Example

ridge is a bridge to convert request/response payload of API Gateway / ALB / Function URLs to Go's net/http.Request and net/http.ResponseWriter.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/fujiwara/ridge"
)

func main() {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)
	ridge.Run(":8080", "/", mux)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello %s\n", r.FormValue("name"))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Hello World")
	fmt.Fprintln(w, r.URL)
}
```

This code works on AWS Lambda and also as a standalone HTTP server.

If you run this code on AWS Lambda, you can access the function via API Gateway,  ALB,  or Function URLs. ridge converts the request payload to Go's net/http.Request, and the response that your app returns converts to the payload of Lambda.

You can access the function via Go's net/http server if you run this code as a standalone HTTP server.

So you can develop and test your Lambda function locally without deploying to AWS.

Also, you can switch the runtime environment between AWS Lambda and a standalone HTTP server on Amazon ECS, Amazon EC2, or AWS AppRunner **without any code and built binary changes**.

### ridge.Run(address, prefix, handler)

`ridge.Run(address, prefix, handler)` works as below.

- If a process is running on Lambda (`AWS_EXECUTION_ENV` or `AWS_LAMBDA_RUNTIME_API` environment variable defined),
  - Call lambda.Start()
- Otherwise start a net/http server using path prefix and address.
  - path prefix is used to strip the prefix from the request path.

### ridge.ToRequestV1(*http.Request)

`ridge.ToRequestV1(*http.Request)` converts a net/http.Request to an API Gateway V1 event payload.

This function helps call your ridge application handler directly from another service that can invoke Lambda.

For example, you can invoke ridge application on Lambda from EventBridge, Step Functions, or another Lambda function with the payload that `ridge.ToRequestV1` returns.

```go
req, _ := http.NewRequest("GET", "http://example.com/hello?name=ridge", nil)
payload, _ := ridge.ToRequestV1(req)
b, _ := json.Marshal(payload)
input := &lambda.InvokeInput{
	FunctionName:   aws.String("your-ridge-function"),
	InvocationType: types.InvocationTypeRequestResponse,
	Payload:        b,
}
result, _ := svc.Invoke(input)
// ...
```

### ridge.ToRequestV2(*http.Request)

`ridge.ToRequestV2(*http.Request)` converts a net/http.Request to an API Gateway V2 event payload.

### Custom request builder

You can use a custom request builder to convert the AWS Lambda invoke payload to net/http.Request.

```go
r := ridge.New(":8080", "/", mux)
r.RequestBuilder = func(payload json.RawMessage) (*http.Request, error) {
    // your custom request builder
}
r.RunWithContext(ctx)
```

default request builder is `ridge.NewRequest`.

### SIGTERM handler

ridge catches SIGTERM and stops the server gracefully by default.

You can add a custom handler to the SIGTERM signal.

```go
r := ridge.New(":8080", "/", mux)
r.SIGTERMHandler = func() {
	// your custom handler
}
r.RunWithContext(ctx)
```

ridge calls the handler when it receives a SIGTERM signal. After the handler returns, ridge stops the server.

**Note**

When you run ridge on AWS Lambda, ridge uses [`lambda.WithEnableSIGTERM`](https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#WithEnableSIGTERM) to call the handler.

You must add [Lambda external extensions](https://docs.aws.amazon.com/lambda/latest/dg/lambda-extensions.html) at least one to the handler. If without extensions, the handler will not be called because the Lambda runtime does not send the SIGTERM signal.

The handler must return in 500ms.

## LICENSE

The MIT License (MIT)

Copyright (c) 2016- FUJIWARA Shunichiro
