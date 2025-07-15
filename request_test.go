package ridge_test

import (
	"encoding/json"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestInvalidRequest(t *testing.T) {
	invalidPayload := json.RawMessage(`{"foo":"bar"}`)
	_, _, err := ridge.NewRequest(invalidPayload)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	t.Log(err)
}
