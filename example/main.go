package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fujiwara/ridge"
)

var mux = http.NewServeMux()

func init() {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)
}

func main() {
	ridge.ProxyProtocol = true
	ridge.Run(":8080", "/", mux)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	log.Println("handleHello", r.URL)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello %s from %s\n", r.FormValue("name"), r.RemoteAddr)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("handleRoot", r.URL)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello World from %s", r.RemoteAddr)
	fmt.Fprintln(w, r.URL)
}
