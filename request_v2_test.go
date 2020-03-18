package ridge_test

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestGetRequestV2(t *testing.T) {
	f, err := os.Open("test/get-v2.json")
	if err != nil {
		t.Fatalf("failed to open test/get-v2.json: %s", err)
	}
	body, _ := ioutil.ReadAll(f)
	r, err := ridge.NewRequest(json.RawMessage(body))
	if err != nil {
		t.Fatalf("failed to decode getEvent: %s", err)
	}
	if r.Host != "abcdefg.execute-api.ap-northeast-1.amazonaws.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "GET" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/?foo=1&foo=2&bar=3")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.FormValue("foo"); v != "1" {
		t.Errorf("FormValue(foo): %s is not expected", v)
	}
	if v := r.Form["foo"][0]; v != "1" {
		t.Errorf("FormValue(foo(0)): %s is not expected", v)
	}
	if v := r.Form["foo"][1]; v != "2" {
		t.Errorf("FormValue(foo(1)): %s is not expected", v)
	}
	if v := r.Header.Get("x-amzn-trace-id"); v != "Root=1-5e723e47-2ffda90008b1b60064fac400" {
		t.Errorf("Header[x-amzn-trace-id]: %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
}

func TestPostRequestV2(t *testing.T) {
	f, err := os.Open("test/post-v2.json")
	if err != nil {
		t.Fatalf("failed to open test/post-v2.json: %s", err)
	}
	body, _ := ioutil.ReadAll(f)
	r, err := ridge.NewRequest(json.RawMessage(body))
	if err != nil {
		t.Fatalf("failed to decode postEvent: %s", err)
	}

	if r.Host != "abcdefg.execute-api.ap-northeast-1.amazonaws.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "POST" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.FormValue("foo"); v != "bar baz" {
		t.Errorf("PostFormValue(foo): %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
	if v := r.Header.Get("x-amzn-trace-id"); v != "Root=1-5e723db7-6077c85e0d781094f0c83e24" {
		t.Errorf("Header[x-amzn-trace-id]: %s is not expected", v)
	}
	if r.ContentLength != 13 {
		t.Errorf("Content-Length: %d is not expected", r.ContentLength)
	}
}
