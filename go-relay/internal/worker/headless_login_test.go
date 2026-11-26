package worker

import (
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
