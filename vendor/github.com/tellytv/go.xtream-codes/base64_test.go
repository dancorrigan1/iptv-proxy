package xtreamcodes

import (
	"bytes"
	"testing"
)

func TestNewFromStringWithPaddingRoundTrip(t *testing.T) {
	encoded := "dGVzdA=="

	value, err := NewFromString(encoded)
	if err != nil {
		t.Fatalf("NewFromString returned error: %v", err)
	}

	if len(*value) == 0 {
		t.Fatalf("expected decoded value to be non-empty")
	}

	marshaled, err := value.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}

	expectedJSON := "\"" + encoded + "\""
	if string(marshaled) != expectedJSON {
		t.Fatalf("expected JSON %s, got %s", expectedJSON, marshaled)
	}

	var decoded Base64Value
	if err := decoded.UnmarshalJSON(marshaled); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}

	if len(decoded) == 0 {
		t.Fatalf("expected unmarshaled value to be non-empty")
	}

	if !bytes.Equal([]byte(decoded), []byte(*value)) {
		t.Fatalf("expected round-trip bytes to match")
	}
}

func TestNewFromStringWithPlusSlashRoundTrip(t *testing.T) {
	encoded := "+/8="

	value, err := NewFromString(encoded)
	if err != nil {
		t.Fatalf("NewFromString returned error: %v", err)
	}

	if len(*value) == 0 {
		t.Fatalf("expected decoded value to be non-empty")
	}

	marshaled, err := value.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}

	expectedJSON := "\"" + encoded + "\""
	if string(marshaled) != expectedJSON {
		t.Fatalf("expected JSON %s, got %s", expectedJSON, marshaled)
	}

	var decoded Base64Value
	if err := decoded.UnmarshalJSON(marshaled); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}

	if len(decoded) == 0 {
		t.Fatalf("expected unmarshaled value to be non-empty")
	}

	if !bytes.Equal([]byte(decoded), []byte(*value)) {
		t.Fatalf("expected round-trip bytes to match")
	}
}
