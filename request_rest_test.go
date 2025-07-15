package ridge_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestRESTAPIRequestProcessing(t *testing.T) {
	payload, err := os.ReadFile("test/get-rest.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	req, _, err := ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	// Verify basic request properties
	if req.Method != "GET" {
		t.Errorf("expected method GET, got %s", req.Method)
	}

	if req.URL.Path != "/path/to/example" {
		t.Errorf("expected path /path/to/example, got %s", req.URL.Path)
	}

	if req.Header.Get("User-Agent") != "curl/7.43.0" {
		t.Errorf("expected User-Agent curl/7.43.0, got %s", req.Header.Get("User-Agent"))
	}

	// Verify query parameters
	if req.URL.Query().Get("foo") != "bar baz" {
		t.Errorf("expected query param foo=bar baz, got %s", req.URL.Query().Get("foo"))
	}

	// Verify Host header
	if req.Host != "abcdefg.execute-api.ap-northeast-1.example.com" {
		t.Errorf("expected host abcdefg.execute-api.ap-northeast-1.example.com, got %s", req.Host)
	}

	// Verify request ID header is set
	if req.Header.Get("X-Amzn-RequestId") != "817175b9-890f-11e6-960e-4f321627a748" {
		t.Errorf("expected request ID header, got %s", req.Header.Get("X-Amzn-RequestId"))
	}
}

func TestPayloadVersionForced(t *testing.T) {
	// Test the PayloadVersion global variable functionality
	originalVersion := ridge.PayloadVersion
	defer func() {
		ridge.PayloadVersion = originalVersion
	}()

	payload, err := os.ReadFile("test/get-rest.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// Force version to 2.0
	ridge.PayloadVersion = "2.0"

	req, _, err := ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		// This is expected to fail because REST API payload structure doesn't match v2.0 format
		t.Logf("Expected failure when forcing v2.0 on REST API payload: %v", err)
	} else if req != nil {
		t.Log("Request processed successfully with forced version")
	}

	// Reset to empty to test auto-detection
	ridge.PayloadVersion = ""

	req2, _, err := ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	if req2 == nil {
		t.Fatal("NewRequest returned nil request")
	}
}

func TestUnsupportedPayloadVersion(t *testing.T) {
	// Create a payload with an unsupported version
	unsupportedPayload := `{"version": "3.0", "path": "/test"}`

	_, _, err := ridge.NewRequest(json.RawMessage(unsupportedPayload))
	if err == nil {
		t.Fatal("NewRequest should have failed for unsupported version")
	}

	expectedError := "payload Version 3.0 is not supported"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
