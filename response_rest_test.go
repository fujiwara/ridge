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

	// Create ResponseWriter
	w := ridge.NewResponseWriter()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Set-Cookie", "session=abc123")
	w.Header().Add("Set-Cookie", "token=xyz789")
	w.WriteHeader(200)
	w.WriteString(`{"message": "hello"}`)

	// Use ResponseFor with empty version for REST API
	resp := w.ResponseFor("")

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

	req, err := ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	// Verify that HTTP API has version header
	version := req.Header.Get(ridge.PayloadVersionHeaderName)
	if version != "2.0" {
		t.Errorf("expected version header 2.0, got %s", version)
	}

	// Create ResponseWriter
	w := ridge.NewResponseWriter()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Set-Cookie", "session=abc123")
	w.WriteHeader(200)
	w.WriteString(`{"message": "hello"}`)

	// Use ResponseFor with version from request
	resp := w.ResponseFor(version)

	// HTTP API responses should have cookies field
	if len(resp.Cookies) != 1 {
		t.Errorf("HTTP API response should have cookies field, got %d cookies", len(resp.Cookies))
	}

	if resp.Cookies[0] != "session=abc123" {
		t.Errorf("unexpected cookie value: %s", resp.Cookies[0])
	}
}
