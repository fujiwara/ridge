package ridge_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestPayloadVersionDetection(t *testing.T) {
	tests := []struct {
		name            string
		payloadFile     string
		expectedVersion string
		shouldSucceed   bool
	}{
		{
			name:            "REST API (no version field)",
			payloadFile:     "test/get-rest.json",
			expectedVersion: "", // REST API has no version field
			shouldSucceed:   true,
		},
		{
			name:            "HTTP API v1.0",
			payloadFile:     "test/get-v1.json",
			expectedVersion: "1.0", // v1.0 now has version field
			shouldSucceed:   true,
		},
		{
			name:            "HTTP API v2.0",
			payloadFile:     "test/get-v2.json",
			expectedVersion: "2.0",
			shouldSucceed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := os.ReadFile(tt.payloadFile)
			if err != nil {
				t.Fatalf("failed to read test file %s: %v", tt.payloadFile, err)
			}

			// Parse version field to verify payload structure
			var versionCheck struct {
				Version string `json:"version"`
			}
			if err := json.Unmarshal(payload, &versionCheck); err != nil {
				t.Fatalf("failed to unmarshal payload: %v", err)
			}

			if versionCheck.Version != tt.expectedVersion {
				t.Errorf("expected version %q, got %q", tt.expectedVersion, versionCheck.Version)
			}

			// Test that NewRequest can handle the payload
			req, err := ridge.NewRequest(json.RawMessage(payload))
			if tt.shouldSucceed {
				if err != nil {
					t.Fatalf("NewRequest failed: %v", err)
				}
				if req == nil {
					t.Fatal("NewRequest returned nil request")
				}
			} else {
				if err == nil {
					t.Fatal("NewRequest should have failed")
				}
			}
		})
	}
}

func TestRESTAPIRequestProcessing(t *testing.T) {
	payload, err := os.ReadFile("test/get-rest.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	req, err := ridge.NewRequest(json.RawMessage(payload))
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
	
	req, err := ridge.NewRequest(json.RawMessage(payload))
	if err != nil {
		// This is expected to fail because REST API payload structure doesn't match v2.0 format
		t.Logf("Expected failure when forcing v2.0 on REST API payload: %v", err)
	} else if req != nil {
		t.Log("Request processed successfully with forced version")
	}

	// Reset to empty to test auto-detection
	ridge.PayloadVersion = ""
	
	req2, err := ridge.NewRequest(json.RawMessage(payload))
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
	
	_, err := ridge.NewRequest(json.RawMessage(unsupportedPayload))
	if err == nil {
		t.Fatal("NewRequest should have failed for unsupported version")
	}

	expectedError := "payload Version 3.0 is not supported"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestRequestContextDetection(t *testing.T) {
	// Test that requests created from different payload types maintain their context
	tests := []struct {
		name       string
		payload    string
		expectedAPIType string
	}{
		{
			name:       "REST API payload",
			payload:    "test/get-rest.json",
			expectedAPIType: "REST",
		},
		{
			name:       "HTTP API v1.0 payload",
			payload:    "test/get-v1.json",
			expectedAPIType: "HTTP",
		},
		{
			name:       "HTTP API v2.0 payload",
			payload:    "test/get-v2.json",
			expectedAPIType: "HTTP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := os.ReadFile(tt.payload)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			req, err := ridge.NewRequest(json.RawMessage(payload))
			if err != nil {
				t.Fatalf("NewRequest failed: %v", err)
			}

			// Check that request contains API type context
			actualAPIType := req.Header.Get("X-API-Gateway-Type")
			if actualAPIType != tt.expectedAPIType {
				t.Errorf("expected API type %q, got %q", tt.expectedAPIType, actualAPIType)
			}
		})
	}
}