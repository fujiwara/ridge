package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fujiwara/ridge"
)

var mux = http.NewServeMux()

func init() {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)
}

func main() {
	app := ridge.New(":8080", "/", mux)
	app.ProxyProtocol = true
	app.TermHandler = func() {
		log.Println("TERM signal received")
		time.Sleep(100 * time.Millisecond)
		log.Println("Goodbye")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.RunWithContext(ctx)
	log.Println("shutdown complete")
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
