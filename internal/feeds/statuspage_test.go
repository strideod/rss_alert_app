package feeds

import (
	"testing"
)

func TestIncidentKeyFromLink(t *testing.T) {
	tests := []struct {
		link     string
		wantKey  string
		wantBool bool
	}{
		{
			link:     "https://status.example.com/incidents/abc123",
			wantKey:  "abc123",
			wantBool: true,
		},
		{
			link:     "https://status.example.com/incidents/2024/06/01/incident-title",
			wantKey:  "incident-title",
			wantBool: true,
		},
		{
			link:     "https://status.example.com/some/other/path",
			wantKey:  "path",
			wantBool: true,
		},
		{
			link:     "invalid-url",
			wantKey:  "",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		gotKey, gotBool := IncidentKeyFromLink(tt.link)
		if gotKey != tt.wantKey || gotBool != tt.wantBool {
			t.Errorf("IncidentKeyFromLink(%q) = (%q, %v); want (%q, %v)", tt.link, gotKey, gotBool, tt.wantKey, tt.wantBool)
		}
	}
}
// key = 
// status =
// hash stable