package ridge_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestGetRequest(t *testing.T) {
	f, err := os.Open("test/get.json")
	if err != nil {
		t.Fatalf("failed to open test/get.json: %s", err)
	}
	body, _ := ioutil.ReadAll(f)
	r, err := ridge.NewRequest(json.RawMessage(body))
	if err != nil {
		t.Fatalf("failed to decode getEvent: %s", err)
	}
	if r.Host != "abcdefg.execute-api.ap-northeast-1.example.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "GET" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/path/to/example?foo=bar+baz&foo=boo+uoo")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.FormValue("foo"); v != "bar baz" {
		t.Errorf("FormValue(foo): %s is not expected", v)
	}
	if v := r.Form["foo"][0]; v != "bar baz" {
		t.Errorf("FormValue(foo(0)): %s is not expected", v)
	}
	if v := r.Form["foo"][1]; v != "boo uoo" {
		t.Errorf("FormValue(foo(1)): %s is not expected", v)
	}
	if v := r.Header.Get("CloudFront-Viewer-Country"); v != "JP" {
		t.Errorf("Header[CloudFront-Viewer-Country]: %s is not expected", v)
	}
	if v := r.Header.Get("Via"); v != "1.1 a3fed41c60e2fab219a274640e58ebe5.cloudfront.net (CloudFront)" {
		t.Errorf("Header[Via]: %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
}

func TestPostRequest(t *testing.T) {
	f, err := os.Open("test/post.json")
	if err != nil {
		t.Fatalf("failed to open test/post.json: %s", err)
	}
	body, _ := ioutil.ReadAll(f)
	r, err := ridge.NewRequest(json.RawMessage(body))
	if err != nil {
		t.Fatalf("failed to decode postEvent: %s", err)
	}

	if r.Host != "abcdefg.execute-api.ap-northeast-1.example.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "POST" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/path/to/example")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.FormValue("foo"); v != "bar baz" {
		t.Errorf("PostFormValue(foo): %s is not expected", v)
	}
	if v := r.Header.Get("CloudFront-Viewer-Country"); v != "JP" {
		t.Errorf("Header[CloudFront-Viewer-Country]: %s is not expected", v)
	}
	if v := r.Header.Get("Via"); v != "1.1 736a82fbf158fe646f468bd5664ef95c.cloudfront.net (CloudFront)" {
		t.Errorf("Header[Via]: %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
	if r.ContentLength != 13 {
		t.Errorf("Content-Length: %d is not expected", r.ContentLength)
	}
}

func TestBase64EncodedRequest(t *testing.T) {
	f, err := os.Open("test/base64.json")
	if err != nil {
		t.Fatalf("failed to open test/post.json: %s", err)
	}
	body, _ := ioutil.ReadAll(f)
	r, err := ridge.NewRequest(json.RawMessage(body))
	if err != nil {
		t.Fatalf("failed to decode postEvent: %s", err)
	}

	if r.Host != "abcdefg.execute-api.ap-northeast-1.example.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "POST" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/path/to/example")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.FormValue("foo"); v != "bar baz" {
		t.Errorf("PostFormValue(foo): %s is not expected", v)
	}
	if v := r.Header.Get("CloudFront-Viewer-Country"); v != "JP" {
		t.Errorf("Header[CloudFront-Viewer-Country]: %s is not expected", v)
	}
	if v := r.Header.Get("Via"); v != "1.1 736a82fbf158fe646f468bd5664ef95c.cloudfront.net (CloudFront)" {
		t.Errorf("Header[Via]: %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
	if r.ContentLength != 13 {
		t.Errorf("Content-Length: %d is not expected", r.ContentLength)
	}
}

func TestResponseWriter(t *testing.T) {
	w := ridge.NewResponseWriter()

	for _, s := range []string{"abcd", "efgh"} {
		n, err := io.WriteString(w, s)
		if n != 4 {
			t.Error("invalid wrote bytes length", n)
		}
		if err != nil {
			t.Error(err)
		}
	}

	w.WriteHeader(500)
	w.Header().Add("Foo", "foo")
	w.Header().Add("Bar", "bar1")
	w.Header().Add("Bar", "bar2")
	res := w.Response()
	if res.StatusCode != 500 {
		t.Error("unexpected status code", res.StatusCode)
	}
	if res.Headers["Foo"] != "foo" {
		t.Error("unexpected Header Foo", res.Headers["Foo"])
	}
	if res.Headers["Bar"] != "bar1" {
		t.Error("unexpected Header Bar", res.Headers["Bar"])
	}
	if res.Body != "abcdefgh" {
		t.Error("unexpected Header Bar", res.Headers["Bar"])
	}
	if res.IsBase64Encoded != false {
		t.Error("set isBase64Encoded = true, but this is text response")
	}
	if res.Headers["Content-Type"] != "text/plain" {
		t.Error("invalid content-type")
	}
}

func TestResponseWriter__Image(t *testing.T) {
	bs, err := ioutil.ReadFile("test/bluebox.png")
	if err != nil {
		t.Error(err)
	}
	expectedBody := base64.StdEncoding.EncodeToString(bs)

	w := ridge.NewResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "http://example.com/bluebox.png", nil)
	if err != nil {
		t.Error(err)
	}
	http.ServeFile(w, req, "test/bluebox.png")

	res := w.Response()
	if res.StatusCode != http.StatusOK {
		t.Errorf("response status is not StatusOK: %d", res.StatusCode)
	}
	if res.IsBase64Encoded != true {
		t.Error("isBase64Encoded is not true")
	}
	if res.Body != expectedBody {
		t.Errorf("base64 encoded body is not match: got=%s", res.Body)
	}
	if res.Headers["Content-Type"] != "image/png" {
		t.Error("invalid content-type")
	}
}
