package xtreamcodes

// Package base64 provides a byte slice type that marshals into json as
// a raw (no padding) base64url value.
// Originally from https://github.com/manifoldco/go-base64

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
	"sync"
)

// Base64Value is a base64url encoded json object,
type Base64Value []byte

var (
	base64EncodingCache = struct {
		sync.RWMutex
		data map[*Base64Value]*base64.Encoding
	}{data: make(map[*Base64Value]*base64.Encoding)}
)

// New returns a pointer to a Base64Value, cast from the given byte slice.
// This is a convenience function, handling the address creation that a direct
// cast would not allow.
func New(b []byte) *Base64Value {
	v := Base64Value(b)
	setBase64Encoding(&v, base64.RawURLEncoding)
	return &v
}

// NewFromString returns a Base64Value containing the decoded data in encoded.
func NewFromString(encoded string) (*Base64Value, error) {
	out, enc, err := decodeBase64String(encoded)
	if err != nil {
		return nil, err
	}

	v := Base64Value(out)
	setBase64Encoding(&v, enc)
	return &v, nil
}

// MarshalJSON returns the ba64url encoding of bv for JSON representation.
func (bv *Base64Value) MarshalJSON() ([]byte, error) {
	return []byte("\"" + encodeBase64String(bv) + "\""), nil
}

func (bv *Base64Value) String() string {
	return encodeBase64String(bv)
}

// UnmarshalJSON sets bv to the bytes represented in the base64url encoding b.
func (bv *Base64Value) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return errors.New("value is not a string")
	}

	decoded, enc, err := decodeBase64String(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	v := reflect.ValueOf(bv).Elem()
	v.SetBytes(decoded)
	setBase64Encoding(bv, enc)
	return nil
}

func encodeBase64String(bv *Base64Value) string {
	if bv == nil {
		return ""
	}

	enc := getBase64Encoding(bv)
	return enc.EncodeToString(*bv)
}

func decodeBase64String(encoded string) ([]byte, *base64.Encoding, error) {
	candidates := candidateEncodings(encoded)
	var lastErr error

	for _, enc := range candidates {
		if out, err := enc.DecodeString(encoded); err == nil {
			return out, enc, nil
		} else {
			lastErr = err
		}
	}

	if lastErr == nil {
		lastErr = errors.New("invalid base64 value")
	}

	return nil, nil, lastErr
}

func candidateEncodings(encoded string) []*base64.Encoding {
	guess := guessEncoding(encoded)
	encs := []*base64.Encoding{guess}

	for _, enc := range []*base64.Encoding{
		base64.StdEncoding,
		base64.URLEncoding,
		base64.RawStdEncoding,
		base64.RawURLEncoding,
	} {
		if enc != guess {
			encs = append(encs, enc)
		}
	}

	return encs
}

func guessEncoding(encoded string) *base64.Encoding {
	padded := strings.ContainsRune(encoded, '=')

	for i := 0; i < len(encoded); i++ {
		switch encoded[i] {
		case '+', '/':
			if padded {
				return base64.StdEncoding
			}
			return base64.RawStdEncoding
		case '-', '_':
			if padded {
				return base64.URLEncoding
			}
			return base64.RawURLEncoding
		}
	}

	if padded {
		return base64.StdEncoding
	}

	return base64.RawURLEncoding
}

func setBase64Encoding(bv *Base64Value, enc *base64.Encoding) {
	if bv == nil {
		return
	}

	base64EncodingCache.Lock()
	base64EncodingCache.data[bv] = enc
	base64EncodingCache.Unlock()
}

func getBase64Encoding(bv *Base64Value) *base64.Encoding {
	if bv == nil {
		return base64.RawURLEncoding
	}

	base64EncodingCache.RLock()
	enc, ok := base64EncodingCache.data[bv]
	base64EncodingCache.RUnlock()
	if ok && enc != nil {
		return enc
	}

	return base64.RawURLEncoding
}
