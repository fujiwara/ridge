package ridge_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/fujiwara/ridge"
	"github.com/google/go-cmp/cmp"
)

func TestGetRequest(t *testing.T) {
	f, err := os.Open("test/get-v1.json")
	if err != nil {
		t.Fatalf("failed to open test/get-v1.json: %s", err)
	}
	body, _ := io.ReadAll(f)
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
	if v := r.Header.Get("X-Amzn-RequestId"); v != "817175b9-890f-11e6-960e-4f321627a748" {
		t.Errorf("Header[X-Amzn-RequestId]: %s is not expected", v)
	}
	// Verify version header is set for v1.0
	if v := r.Header.Get(ridge.PayloadVersionHeaderName); v != "1.0" {
		t.Errorf("expected version header 1.0, got %s", v)
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
	body, _ := io.ReadAll(f)
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
	if v := r.Header.Get("X-Amzn-RequestId"); v != "8eed9b4f-890f-11e6-9f3c-1584342606cd" {
		t.Errorf("Header[X-Amzn-RequestId]: %s is not expected", v)
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
	body, _ := io.ReadAll(f)
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
	if v := r.Header.Get("X-Amzn-RequestId"); v != "8eed9b4f-890f-11e6-9f3c-1584342606cd" {
		t.Errorf("Header[X-Amzn-RequestId]: %s is not expected", v)
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
	w.Header().Add("Set-Cookie", "cookie1=value1; Secure; HttpOnly")
	w.Header().Add("Set-Cookie", "cookie2=value2; Domain=example.com; Path=/; Max-Age=3600; Secure; HttpOnly")
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
	if res.Headers["Content-Type"] != "text/plain; charset=utf-8" {
		t.Error("invalid content-type")
	}
	if res.Cookies[0] != "cookie1=value1; Secure; HttpOnly" {
		t.Error("invalid cookie", res.Cookies[0])
	}
	if res.Cookies[1] != "cookie2=value2; Domain=example.com; Path=/; Max-Age=3600; Secure; HttpOnly" {
		t.Error("invalid cookie", res.Cookies[1])
	}
}

func TestResponseWriter__Image(t *testing.T) {
	bs, err := os.ReadFile("test/bluebox.png")
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

var roundTripRequest = map[string]func() *http.Request{
	"GET": func() *http.Request {
		r, _ := http.NewRequest(http.MethodGet, "https://example.com/path/to/example?foo=bar+baz&foo=boo+uoo", nil)
		r.Header.Set("Cookie", "foo=bar; xxx=yyy")
		r.Header.Set("x-api-key", "zzz")
		return r
	},
	"POST": func() *http.Request {
		r, _ := http.NewRequest(http.MethodPost, "https://example.com/path/to/example", nil)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("x-api-key", "zzz")
		body := `{"foo":"bar baz"}`
		r.Body = io.NopCloser(strings.NewReader(body))
		r.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))
		return r
	},
	"MultiValue": func() *http.Request {
		r, _ := http.NewRequest(http.MethodGet, "https://example.com/path/to/example?foo=1&foo=2", nil)
		r.Header.Add("foo", "bar")
		r.Header.Add("foo", "baz")
		return r
	},
}

func TestV1RoundTrip(t *testing.T) {
	for name, newRequest := range roundTripRequest {
		t.Run(name, func(t *testing.T) {
			or := newRequest()
			od, _ := httputil.DumpRequest(or, true)

			payload, err := ridge.ToRequestV1(or)
			if err != nil {
				t.Error(err)
				return
			}
			b, _ := json.Marshal(payload)
			rr, err := ridge.NewRequest(json.RawMessage(b))
			if err != nil {
				t.Error("failed to decode RequestV1", err)
			}

			// Remove X-Ridge-Payload-Version header before comparison
			// as it's added by ridge for internal use
			rr.Header.Del(ridge.PayloadVersionHeaderName)
			rd, _ := httputil.DumpRequest(rr, true)
			t.Logf("original request: %s", od)
			t.Logf("decoded request: %s", rd)
			if d := cmp.Diff(od, rd); d != "" {
				t.Error("request is not match", d)
			}
		})
	}
}

func TestMinimalValidRequestV1(t *testing.T) {
	payload := json.RawMessage(`{"path": "/path/to/example","httpMethod": "GET"}`)
	r, err := ridge.NewRequest(payload)
	if err != nil {
		t.Error("failed to decode minimal valid RequestV1", err)
	}
	if r.Method != "GET" {
		t.Error("unexpected method", r.Method)
	}
	if r.URL.Path != "/path/to/example" {
		t.Error("unexpected path", r.URL.Path)
	}
}
