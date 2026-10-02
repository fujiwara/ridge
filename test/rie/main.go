// Command rie is a test application for the integration tests with
// the AWS Lambda Runtime Interface Emulator (see rie_test.go).
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fujiwara/ridge"
)

// TermMessage is logged when the TermHandler is called.
const TermMessage = "ridge-rie: TERM signal received"

// Echo is the response body of the "/echo" handler.
type Echo struct {
	Method     string              `json:"method"`
	RequestURI string              `json:"request_uri"`
	RemoteAddr string              `json:"remote_addr"`
	Header     http.Header         `json:"header"`
	Form       map[string][]string `json:"form"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "foo", Value: "1"})
		http.SetCookie(w, &http.Cookie{Name: "bar", Value: "2"})
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/binary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte{0x00, 0x01, 0x02, 0xff})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Echo{
			Method:     r.Method,
			RequestURI: r.URL.RequestURI(),
			RemoteAddr: r.RemoteAddr,
			Header:     r.Header,
			Form:       r.Form,
		})
	})
	app := ridge.New(":8080", "/", mux)
	app.TermHandler = func() {
		log.Println(TermMessage)
	}
	app.Run()
}
