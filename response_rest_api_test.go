package ridge_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestResponseForRESTAPI(t *testing.T) {
	// Load REST API payload and create request
	payload, err := os.ReadFile("test/get-rest.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	_, err = ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	// Create ResponseWriter with REST API type
	w := ridge.NewResponseWriter("REST")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Set-Cookie", "session=abc123")
	w.Header().Add("Set-Cookie", "token=xyz789")
	w.WriteHeader(200)
	w.WriteString(`{"message": "hello"}`)

	resp := w.Response()

	// REST API responses should not have cookies field
	if len(resp.Cookies) != 0 {
		t.Errorf("REST API response should not have cookies field, got %d cookies", len(resp.Cookies))
	}

	// But Set-Cookie should still be in MultiValueHeaders for proper HTTP handling
	cookies := resp.MultiValueHeaders["Set-Cookie"]
	if len(cookies) != 2 {
		t.Errorf("expected 2 Set-Cookie headers in MultiValueHeaders, got %d", len(cookies))
	}
	if cookies[0] != "session=abc123" {
		t.Errorf("unexpected first cookie in MultiValueHeaders: %s", cookies[0])
	}
	if cookies[1] != "token=xyz789" {
		t.Errorf("unexpected second cookie in MultiValueHeaders: %s", cookies[1])
	}

	// Other response fields should work normally
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", resp.Headers["Content-Type"])
	}
	if resp.Body != `{"message": "hello"}` {
		t.Errorf("unexpected body: %s", resp.Body)
	}
}

func TestResponseForHTTPAPI(t *testing.T) {
	// Load HTTP API v2.0 payload and create request
	payload, err := os.ReadFile("test/get-v2.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	_, err = ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	// Create ResponseWriter with HTTP API type
	w := ridge.NewResponseWriter("HTTP")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Set-Cookie", "session=abc123")
	w.WriteHeader(200)
	w.WriteString(`{"message": "hello"}`)

	resp := w.Response()

	// HTTP API responses should have cookies field
	if len(resp.Cookies) != 1 {
		t.Errorf("HTTP API response should have cookies field, got %d cookies", len(resp.Cookies))
	}

	if resp.Cookies[0] != "session=abc123" {
		t.Errorf("unexpected cookie value: %s", resp.Cookies[0])
	}
}

var testRESTAPIResponses = []struct {
	Name string
	JSON []byte
}{
	{
		Name: "rest_api_plain",
		// Note: REST API response doesn't have "cookies" field
		JSON: []byte(`{"statusCode":200,"headers":{"Content-Type":"text/plain"},"multiValueHeaders":{"Content-Type":["text/plain"],"Set-Cookie":["foo=bar","bar=baz"]},"body":"Hello REST","isBase64Encoded":false}`),
	},
}

func TestRESTAPIResponseFormat(t *testing.T) {
	for _, c := range testRESTAPIResponses {
		t.Run(c.Name, func(t *testing.T) {
			var resp map[string]interface{}
			if err := json.Unmarshal(c.JSON, &resp); err != nil {
				t.Error(err)
			}

			// REST API response should not have cookies field
			if _, exists := resp["cookies"]; exists {
				t.Error("REST API response should not contain 'cookies' field")
			}

			// But should have Set-Cookie in multiValueHeaders
			multiValueHeaders := resp["multiValueHeaders"].(map[string]interface{})
			if setCookies, exists := multiValueHeaders["Set-Cookie"]; !exists {
				t.Error("Set-Cookie should be in multiValueHeaders")
			} else {
				cookies := setCookies.([]interface{})
				if len(cookies) != 2 {
					t.Errorf("expected 2 cookies in multiValueHeaders, got %d", len(cookies))
				}
			}
		})
	}
}

func TestResponseWithoutContext(t *testing.T) {
	// Test behavior with unspecified API type
	w := ridge.NewResponseWriter("")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Set-Cookie", "session=abc123")
	w.WriteHeader(200)
	w.WriteString(`{"message": "hello"}`)

	resp := w.Response()

	// Without context, should default to including cookies
	if len(resp.Cookies) != 1 {
		t.Errorf("Default behavior should include cookies, got %d", len(resp.Cookies))
	}

	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("unexpected Content-Type: %s", resp.Headers["Content-Type"])
	}

	if resp.Body != `{"message": "hello"}` {
		t.Errorf("unexpected body: %s", resp.Body)
	}
}
