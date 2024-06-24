package ridge_test

import (
	"testing"

	"github.com/fujiwara/ridge"
)

func TestIsOnLambdaRuntime(t *testing.T) {
	tests := []struct {
		name                string
		awsExecutionEnv     string
		awsLambdaRuntimeAPI string
		handler             string
		expected            bool
	}{
		{"On AWS Lambda go runtime", "AWS_Lambda_go1.x", "", "handler", true},
		{"On AWS Lambda custom runtime", "", "http://localhost:8080", "handler", true},
		{"On AWS Lambda extension with go runtime", "AWS_Lambda_go1.x", "", "", false},
		{"On AWS Lambda extension with custom runtime", "", "http://localhost:8080", "", false},
		{"Not on AWS Lambda", "", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AWS_EXECUTION_ENV", tt.awsExecutionEnv)
			t.Setenv("AWS_LAMBDA_RUNTIME_API", tt.awsLambdaRuntimeAPI)
			t.Setenv("_HANDLER", tt.handler)

			if got := ridge.IsOnLambdaRuntime(); got != tt.expected {
				t.Errorf("IsOnLambdaRuntime() = %v; want %v", got, tt.expected)
			}
		})
	}
}
