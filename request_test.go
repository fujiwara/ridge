package ridge_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestNewRequestAPITypeDetection(t *testing.T) {
	tests := []struct {
		name            string
		payloadFile     string
		expectedAPIType string
	}{
		{
			name:            "REST API payload",
			payloadFile:     "test/get-rest.json",
			expectedAPIType: "REST",
		},
		{
			name:            "HTTP API v1.0 payload",
			payloadFile:     "test/get-v1.json",
			expectedAPIType: "HTTP",
		},
		{
			name:            "HTTP API v2.0 payload",
			payloadFile:     "test/get-v2.json",
			expectedAPIType: "HTTP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := os.ReadFile(tt.payloadFile)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			req, apiType, err := ridge.NewRequest(json.RawMessage(payload))
			if err != nil {
				t.Fatalf("NewRequest failed: %v", err)
			}
			if req == nil {
				t.Fatal("NewRequest returned nil request")
			}

			if apiType != tt.expectedAPIType {
				t.Errorf("expected API type %q, got %q", tt.expectedAPIType, apiType)
			}
		})
	}
}

func TestInvalidRequest(t *testing.T) {
	invalidPayload := json.RawMessage(`{"foo":"bar"}`)
	_, _, err := ridge.NewRequest(invalidPayload)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	t.Log(err)
}
