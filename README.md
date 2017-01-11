# ridge

AWS Lambda HTTP Proxy integration event bridge to Go net/http.

## Example

ridge is a bridge to convert API Gateway with Lambda Proxy Integration request/response and net/http.Request and net/http.ResponseWriter.

- API Gateway with Lambda Proxy Integration through a Proxy Resource
- [Apex](http://apex.run/)

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/fujiwara/ridge"
)

var mux = http.NewServeMux()

func init() {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)
}

func main() {
	ridge.Run(":8080", "/api", mux)
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

1. Run `apex init`.
2. Place main.go to functions/example/.
3. Edit project.json
  - `"language": "go"`
4. Run `apex deploy`
5. Create API Gateway Proxy Integration. http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html

### ridge.Run(address, prefix, handler)

`ridge.Run(address, prefix, handler)` works as below.

- If a process is running on Apex (`APEX_FUNCTION_NAME` environment variable defined), call apex.HandleFunc().
- Otherwise start a net/http server using prefix and address.

## LICENSE

The MIT License (MIT)

Copyright (c) 2016 FUJIWARA Shunichiro
