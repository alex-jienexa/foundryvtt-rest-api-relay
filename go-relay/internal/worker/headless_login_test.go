package worker

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateFoundryJoinResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    string
	}{
		{
			name:       "Foundry v14 native success payload",
			statusCode: 200,
			body:       `{"message":"JOIN.Success","redirect":"/game"}`,
		},
		{
			name:       "legacy explicit success payload",
			statusCode: 200,
			body:       `{"status":"success","redirect":"/game"}`,
		},
		{
			name:       "alternate explicit success payload",
			statusCode: 200,
			body:       `{"response":"success","redirect":"/game"}`,
		},
		{
			name:       "explicit application error",
			statusCode: 200,
			body:       `{"status":"error","message":"Invalid password"}`,
			wantErr:    `status="error"`,
		},
		{
			name:       "Foundry HTTP error",
			statusCode: 401,
			body:       `{"error":"Invalid password"}`,
			wantErr:    "HTTP 401: Invalid password",
		},
		{
			name:       "malformed successful response",
			statusCode: 200,
			body:       `not-json`,
			wantErr:    "invalid JSON response",
		},
		{
			name:       "empty successful response",
			statusCode: 200,
			body:       `null`,
			wantErr:    "empty join response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFoundryJoinResponse(tt.statusCode, []byte(tt.body))
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected success, got error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestFoundryV14JoinPayload(t *testing.T) {
	encoded, err := buildFoundryV14JoinPayload("Puffy", "gkYizCRHVPCjt0Xe", "test-password")
	if err != nil {
		t.Fatalf("build join payload: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(encoded), &got); err != nil {
		t.Fatalf("unmarshal join payload: %v", err)
	}
	want := map[string]string{
		"username": "Puffy",
		"userId":   "gkYizCRHVPCjt0Xe",
		"password": "test-password",
		"action":   "join",
	}
	if len(got) != len(want) {
		t.Fatalf("unexpected payload fields: %#v", got)
	}
	for key, value := range want {
		if got[key] != value {
			t.Errorf("payload[%q] = %q, want %q", key, got[key], value)
		}
	}
}

func TestFoundryV14JoinPayloadRequiresUserID(t *testing.T) {
	if _, err := buildFoundryV14JoinPayload("Puffy", "", "test-password"); err == nil {
		t.Fatal("expected missing FOUNDRY_USER_ID to fail")
	}
}
