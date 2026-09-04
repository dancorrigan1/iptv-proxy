package xtreamproxy

import (
	"encoding/json"
	"testing"

	xtream "github.com/tellytv/go.xtream-codes"
)

func TestEPGInfoUnmarshalStandardBase64(t *testing.T) {
	data := []byte(`{
		"id": "196967674",
		"epg_id": "57",
		"title": "MjAyNiBOQkEgRmluYWxz",
		"lang": "en",
		"start": "2026-06-11 00:30:00",
		"end": "1781146800",
		"description": "VGhlIEtuaWNrcyBob3N0IHRoZSBTcHVycyBhdCBNYWRpc29uIFNxdWFyZSBHYXJkZW4gZm9yIEdhbWUgNCBvZiB0aGUgTkJBIEZpbmFscy4=",
		"channel_id": "abc13ktrktv.us",
		"start_timestamp": "1781137800",
		"stop_timestamp": "1781146800",
		"stop": "2026-06-11 03:00:00"
	}`)

	var epg xtream.EPGInfo
	if err := json.Unmarshal(data, &epg); err != nil {
		t.Fatalf("unmarshal EPGInfo: %v", err)
	}

	if got, want := string(epg.Title), "2026 NBA Finals"; got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}

	wantDescription := "The Knicks host the Spurs at Madison Square Garden for Game 4 of the NBA Finals."
	if got := string(epg.Description); got != wantDescription {
		t.Fatalf("description = %q, want %q", got, wantDescription)
	}
}
