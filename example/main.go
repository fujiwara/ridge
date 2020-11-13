package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fujiwara/ridge"
)

var mux = http.NewServeMux()

func init() {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello", handleHello)
}

func main() {
	ridge.ProxyProtocol = true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		s := <-sigchan
		log.Println("got signal", s)
		cancel()
	}()
	ridge.RunWithContext(ctx, ":8080", "/", mux)
	log.Println("shutdown complate")
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
