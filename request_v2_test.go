package ridge_test

import (
	"encoding/json"
	"io"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/fujiwara/ridge"
	"github.com/google/go-cmp/cmp"
)

func TestGetRequestV2(t *testing.T) {
	f, err := os.Open("test/get-v2.json")
	if err != nil {
		t.Fatalf("failed to open test/get-v2.json: %s", err)
	}
	body, _ := io.ReadAll(f)
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
	u, _ := url.Parse("/%E3%83%9E%E3%83%AB%E3%83%81%E3%83%90%E3%82%A4%E3%83%88?foo=1&foo=2&bar=3")
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
	if v := r.Header.Get("x-amzn-requestid"); v != "Jl6rIhtwNjMEJLQ=" {
		t.Errorf("Header[x-amzn-requestid]: %s is not expected", v)
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
	body, _ := io.ReadAll(f)
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
	if v := r.Header.Get("x-amzn-requestid"); v != "Jl6UpgU9tjMEPLA=" {
		t.Errorf("Header[x-amzn-requestid]: %s is not expected", v)
	}
	if r.ContentLength != 13 {
		t.Errorf("Content-Length: %d is not expected", r.ContentLength)
	}
}

func TestV2RoundTrip(t *testing.T) {
	for name, newRequest := range roundTripRequest {
		t.Run(name, func(t *testing.T) {
			or := newRequest()
			for k, v := range or.Header {
				if len(v) > 1 && k != "Cookie" {
					or.Header.Set(k, strings.Join(v, ",")) // join multiple headers with comma
				}
			}
			od, _ := httputil.DumpRequest(or, true)

			payload, err := ridge.ToRequestV2(or)
			if err != nil {
				t.Error(err)
				return
			}
			b, _ := json.Marshal(payload)
			t.Logf("payload: %s", string(b))
			rr, err := ridge.NewRequest(json.RawMessage(b))
			if err != nil {
				t.Error("failed to decode RequestV1", err)
			}

			rd, _ := httputil.DumpRequest(rr, true)
			t.Logf("original request: %s", od)
			t.Logf("decoded request: %s", rd)
			if d := cmp.Diff(od, rd); d != "" {
				t.Error("request is not match", d)
			}
		})
	}
}

func TestMinimalValidRequestV2(t *testing.T) {
	payload := json.RawMessage(`{
  "version": "2.0",
  "rawPath": "/path/to/example",
  "requestContext": {
    "http": {
      "method": "GET"
    }
  }
}`)
	r, err := ridge.NewRequest(payload)
	if err != nil {
		t.Error("failed to decode minimal valid RequestV2", err)
	}
	if r.Method != "GET" {
		t.Error("unexpected method", r.Method)
	}
	if r.URL.Path != "/path/to/example" {
		t.Error("unexpected path", r.URL.Path)
	}
}
