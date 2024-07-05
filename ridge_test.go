package ridge_test

import (
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
