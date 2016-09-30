# l2h

AWS Lambda HTTP Proxy integration event bridge to Go net/http.

## Example

- API Gateway with Lambda Proxy Integration through a Proxy Resource
- [Apex](http://apex.run/)

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/go-apex"
	"github.com/fujiwara/l2h"
)

var mux = http.NewServeMux()

func init() {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)
}

func main() {
	if os.Getenv("APEX") == "" {
		log.Println("starting up with local httpd")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		r, err := l2h.NewRequest(event)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		w := l2h.NewResponseWriter()
		mux.ServeHTTP(w, r)
		return w.Response(), nil
	})
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

## LICENSE

The MIT License (MIT)

Copyright (c) 2016 FUJIWARA Shunichiro
