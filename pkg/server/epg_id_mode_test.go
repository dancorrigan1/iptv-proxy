package server

import (
	"encoding/json"
	"reflect"
	"testing"

	xtream "github.com/tellytv/go.xtream-codes"
)

func TestRewriteLiveStreamEPGIDsAdvancedResponse(t *testing.T) {
	response := []xtream.Stream{
		{Fields: []byte(`{"stream_id":35481,"epg_channel_id":"msnbc.usa.us","name":"MSNBC HD | USA"}`)},
		{Fields: []byte(`{"stream_id":"1291","epg_channel_id":"teledoce_uruguay_xtream_ai"}`)},
	}

	processed := ProcessResponse(response)
	got := rewriteLiveStreamEPGIDs(processed).([]interface{})
	if got[0].(map[string]interface{})["epg_channel_id"] != "35481" {
		t.Fatalf("first epg_channel_id = %v, want 35481", got[0].(map[string]interface{})["epg_channel_id"])
	}
	if got[1].(map[string]interface{})["epg_channel_id"] != "1291" {
		t.Fatalf("second epg_channel_id = %v, want 1291", got[1].(map[string]interface{})["epg_channel_id"])
	}
	if got[0].(map[string]interface{})["name"] != "MSNBC HD | USA" {
		t.Fatal("unrelated stream fields were changed")
	}
}

func TestRewriteLiveStreamEPGIDsTypedResponse(t *testing.T) {
	response := []xtream.Stream{
		{ID: xtream.FlexInt(35481), EPGChannelID: "msnbc.usa.us", Name: "MSNBC HD | USA"},
	}

	got := rewriteLiveStreamEPGIDs(response).([]xtream.Stream)
	if got[0].EPGChannelID != "35481" {
		t.Fatalf("epg_channel_id = %q, want 35481", got[0].EPGChannelID)
	}
	if response[0].EPGChannelID != "msnbc.usa.us" {
		t.Fatal("input response was mutated")
	}
}

func TestRewriteLiveStreamEPGIDsLeavesMissingOrInvalidIDsAlone(t *testing.T) {
	response := []map[string]interface{}{
		{"epg_channel_id": "missing"},
		{"stream_id": 12.5, "epg_channel_id": "fractional"},
		{"stream_id": json.Number("42"), "epg_channel_id": "number"},
	}

	got := rewriteLiveStreamEPGIDs(response).([]map[string]interface{})
	want := []map[string]interface{}{
		{"epg_channel_id": "missing"},
		{"stream_id": 12.5, "epg_channel_id": "fractional"},
		{"stream_id": json.Number("42"), "epg_channel_id": "42"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("response = %#v, want %#v", got, want)
	}
}
