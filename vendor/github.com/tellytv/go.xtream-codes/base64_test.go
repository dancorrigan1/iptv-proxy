package xtreamcodes

import "testing"

func TestBase64ValueStringMatchesMarshal(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{name: "PaddedStandard", encoded: "dGVzdA=="},
		{name: "URLSafe", encoded: "-_8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value Base64Value
			if err := value.UnmarshalJSON([]byte("\"" + tt.encoded + "\"")); err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if got := value.String(); got != tt.encoded {
				t.Fatalf("String() = %q, want %q", got, tt.encoded)
			}

			marshaled, err := value.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}

			expected := "\"" + tt.encoded + "\""
			if string(marshaled) != expected {
				t.Fatalf("MarshalJSON = %q, want %q", string(marshaled), expected)
			}
		})
	}
}
