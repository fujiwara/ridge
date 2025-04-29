package ridge_test

import (
	"fmt"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestRuntimeEnvironments(t *testing.T) {
	tests := []struct {
		name                string
		awsExecutionEnv     string
		awsLambdaRuntimeAPI string
		handler             string
		onLambdaRuntime     bool
		asLambdaHandler     bool
		asLambdaExtension   bool
	}{
		{"On AWS Lambda go runtime handler", "AWS_Lambda_go1.x", "", "handler", true, true, false},
		{"On AWS Lambda custom runtime handler", "", "http://localhost:8080", "handler", true, true, false},
		{"On AWS Lambda extension with go runtime", "AWS_Lambda_go1.x", "", "", true, false, true},
		{"On AWS Lambda extension with custom runtime", "", "http://localhost:8080", "", true, false, true},
		{"Not on AWS Lambda", "", "", "", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AWS_EXECUTION_ENV", tt.awsExecutionEnv)
			t.Setenv("AWS_LAMBDA_RUNTIME_API", tt.awsLambdaRuntimeAPI)
			t.Setenv("_HANDLER", tt.handler)

			if ridge.OnLambdaRuntime() != tt.onLambdaRuntime {
				t.Errorf("OnLambdaRuntime() = %v, want %v", ridge.OnLambdaRuntime(), tt.onLambdaRuntime)
			}
			if ridge.AsLambdaHandler() != tt.asLambdaHandler {
				t.Errorf("AsLambdaHandler() = %v, want %v", ridge.AsLambdaHandler(), tt.asLambdaHandler)
			}
			if ridge.AsLambdaExtension() != tt.asLambdaExtension {
				t.Errorf("AsLambdaExtension() = %v, want %v", ridge.AsLambdaExtension(), tt.asLambdaExtension)
			}
		})
	}
}

func TestStreamingEnv(t *testing.T) {
	tests := []struct {
		pre   bool
		value string
		want  bool
	}{
		{false, "1", true},
		{false, "0", false},
		{true, "", true},
		{true, "1", true},
		{true, "0", true},
		{false, "true", true},
		{false, "false", false},
		{false, "True", true},
		{false, "False", false},
		{false, "-xxx", false},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("pre=%v,env=%s", tt.pre, tt.value)
		t.Run(name, func(t *testing.T) {
			t.Setenv("RIDGE_STREAMING_RESPONSE", tt.value)
			r := ridge.New(":9999", "/", nil)
			r.StreamingResponse = tt.pre
			r.SetStreamingResponse()
			if r.StreamingResponse != tt.want {
				t.Errorf("StreamingResponse = %v, want %v", r.StreamingResponse, tt.want)
			}
		})
	}
}
