package ridge_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/fujiwara/ridge"
)

var testResponses = []struct {
	Name string
	JSON []byte
}{
	{
		Name: "plain",
		JSON: []byte(`{"statusCode":200,"headers":{"Content-Type":"text/plain"},"multiValueHeaders":{"Content-Type":["text/plain"]},"cookies":["foo=bar","bar=baz"],"body":"Hello XXX","isBase64Encoded":false}`),
	},
	{
		Name: "base64",
		JSON: []byte(`{"statusCode":200,"headers":{"Content-Type":"text/plain"},"multiValueHeaders":{"Content-Type":["text/plain"]},"cookies":["foo=bar","bar=baz"],"body":"SGVsbG8gWFhY","isBase64Encoded":true}`),
	},
}

func TestResponse(t *testing.T) {
	for _, c := range testResponses {
		t.Run(c.Name, func(t *testing.T) {
			var res ridge.Response
			if err := json.Unmarshal(c.JSON, &res); err != nil {
				t.Error(err)
			}
			w := ridge.NewResponseWriter()
			if n, err := res.WriteTo(w); err != nil {
				t.Error(err)
			} else if n != 9 {
				t.Errorf("unexpected body size: %d", n)
			}
			res2 := w.Response()
			if res2.StatusCode != res.StatusCode {
				t.Errorf("unexpected status code: %d", res2.StatusCode)
			}
			if res2.Headers["Content-Type"] != "text/plain" {
				t.Errorf("unexpected Content-Type: %s", res2.Headers["Content-Type"])
			}
			if len(res2.MultiValueHeaders["Content-Type"]) != 1 || res2.MultiValueHeaders["Content-Type"][0] != "text/plain" {
				t.Errorf("unexpected Content-Type: %#v", res2.MultiValueHeaders["Content-Type"])
			}
			if res2.Body != "Hello XXX" {
				t.Errorf("unexpected body: %s", res2.Body)
			}
			if res2.Cookies[0] != "foo=bar" || res2.Cookies[1] != "bar=baz" {
				t.Errorf("unexpected cookies: %#v", res2.Cookies)
			}
			t.Logf("%#v\n", res2)
		})
	}
}

func TestStreamingResponse(t *testing.T) {
	w := ridge.NewStramingResponseWriter()
	signalChan := make(chan struct{}, 1)
	defer close(signalChan)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		defer w.Close()
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		select {
		case <-ctx.Done():
			t.Error("timeout")
			return
		case <-signalChan:
		}
		for i := 0; i < 5; i++ {
			fmt.Fprintf(w, "data: %d\n\n", i)
			w.Flush()
		}
	}()

	w.Wait()
	res := w.Response()
	if res.StatusCode != 200 {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}
	if res.Headers["Content-Type"] != "text/event-stream" {
		t.Errorf("unexpected Content-Type: %s", res.Headers["Content-Type"])
	}
	if len(res.Cookies) != 0 {
		t.Errorf("unexpected cookies: %#v", res.Cookies)
	}

	signalChan <- struct{}{}
	actual, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading body: %v", err)
	}
	if len(actual) == 0 {
		t.Errorf("unexpected empty body")
	}
	expected := `data: 0

data: 1

data: 2

data: 3

data: 4

`
	if string(actual) != expected {
		t.Errorf("unexpected body: %s", string(actual))
	}
}
