package ridge_test

import (
	"encoding/json"
	"testing"

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
