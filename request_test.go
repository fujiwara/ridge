package ridge_test

import (
	"encoding/json"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestInvalidRequest(t *testing.T) {
	invalidPayload := json.RawMessage(`{"foo":"bar"}`)
	_, err := ridge.NewRequest(invalidPayload)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	t.Log(err)
}

func TestRemoteAddr(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected string
	}{
		{"v1 IPv4", `{"version":"1.0","path":"/","httpMethod":"GET","requestContext":{"identity":{"sourceIp":"203.0.113.1"}}}`, "203.0.113.1:0"},
		{"v1 IPv6", `{"version":"1.0","path":"/","httpMethod":"GET","requestContext":{"identity":{"sourceIp":"2001:db8::1"}}}`, "[2001:db8::1]:0"},
		{"v1 without source IP (ALB)", `{"path":"/","httpMethod":"GET"}`, ""},
		{"v2 IPv4", `{"version":"2.0","rawPath":"/","requestContext":{"http":{"method":"GET","sourceIp":"203.0.113.1"}}}`, "203.0.113.1:0"},
		{"v2 IPv6", `{"version":"2.0","rawPath":"/","requestContext":{"http":{"method":"GET","sourceIp":"2001:db8::1"}}}`, "[2001:db8::1]:0"},
		{"v2 without source IP", `{"version":"2.0","rawPath":"/","requestContext":{"http":{"method":"GET"}}}`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ridge.NewRequest(json.RawMessage(tt.payload))
			if err != nil {
				t.Fatal(err)
			}
			if r.RemoteAddr != tt.expected {
				t.Errorf("RemoteAddr = %q, want %q", r.RemoteAddr, tt.expected)
			}
		})
	}
}
