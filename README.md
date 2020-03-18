# ridge

AWS Lambda HTTP Proxy integration event bridge to Go net/http.

## Example

ridge is a bridge to convert API Gateway with Lambda Proxy Integration request/response and net/http.Request and net/http.ResponseWriter.

- API Gateway with Lambda Proxy Integration through a Proxy Resource

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

1. Create IAM role "ridge" for Lambda which have attached policy AWSLambdaBasicExecutionRole.
1. Install [lambroll](https://github.com/fujiwara/lambroll).
1. Place main.go to example/.
1. Run `make deploy` to deploy a lambda function.
1. Create API Gateway with the lambda function.
  - for HTTP API https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-lambda.html
  - for REST API http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html

### ridge.Run(address, prefix, handler)

`ridge.Run(address, prefix, handler)` works as below.

- If a process is running on Lambda (`AWS_EXECUTION_ENV` environment variable defined),
  - Call apex.HandleFunc() when runtime is nodejs*
  - Call lambda.Start() when runtime is go1.x
- Otherwise start a net/http server using prefix and address.

## LICENSE

The MIT License (MIT)

Copyright (c) 2016- FUJIWARA Shunichiro
