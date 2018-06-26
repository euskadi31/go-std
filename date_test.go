package std

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	dateString    = "2012-12-21"
	dateJSON      = []byte(`"` + dateString + `"`)
	nullDateJSON  = []byte(`null`)
	dateValue, _  = time.Parse(dateFormat, dateString)
	badDateObject = []byte(`{"hello": "world"}`)
)

func TestUnmarshalDateJSON(t *testing.T) {
	var ti Date
	err := json.Unmarshal(dateJSON, &ti)
	maybePanic(err)
	assertDate(t, ti, "UnmarshalJSON() json")

	var null Date
	err = json.Unmarshal(nullDateJSON, &null)
	maybePanic(err)
	assertNullDate(t, null, "null time json")

	var invalid Date
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullDate(t, invalid, "invalid from object json")

	var bad Date
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullDate(t, bad, "bad from object json")

	var wrongType Date
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullDate(t, wrongType, "wrong type object json")
}

func TestUnmarshalDateText(t *testing.T) {
	ti := DateFrom(dateValue)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, dateString, "marshal text")

	var unmarshal Date
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertDate(t, unmarshal, "unmarshal text")

	var null Date
	err = null.UnmarshalText(nullDateJSON)
	maybePanic(err)
	assertNullDate(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, string(nullDateJSON), "marshal null text")

	var invalid Date
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDate(t, invalid, "bad string")
}

func TestMarshalDate(t *testing.T) {
	dt := Date{}
	data, err := json.Marshal(dt)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullDateJSON), "null json marshal")

	ti := DateFrom(dateValue)
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(dateJSON), "non-empty json marshal")

	ti.Valid = false
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullDateJSON), "null json marshal")
}

func TestDateFrom(t *testing.T) {
	ti := DateFrom(dateValue)
	assertDate(t, ti, "DateFrom() time.Time")
}

func TestDateFromPtr(t *testing.T) {
	ti := DateFromPtr(&dateValue)
	assertDate(t, ti, "DateFromPtr() time")

	null := DateFromPtr(nil)
	assertNullDate(t, null, "DateFromPtr(nil)")
}

func TestDateSetValid(t *testing.T) {
	var ti time.Time
	change := NewDate(ti, false)
	assertNullDate(t, change, "SetValid()")
	change.SetValid(dateValue)
	assertDate(t, change, "SetValid()")
}

func TestDatePointer(t *testing.T) {
	ti := DateFrom(dateValue)
	ptr := ti.Ptr()
	if *ptr != dateValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, dateValue)
	}

	var nt time.Time
	null := NewDate(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestDateScanValue(t *testing.T) {
	var ti Date
	err := ti.Scan(dateValue)
	maybePanic(err)
	assertDate(t, ti, "scanned time")
	if v, err := ti.Value(); v != dateValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null Date
	err = null.Scan(nil)
	maybePanic(err)
	assertNullDate(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong Date
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDate(t, wrong, "scanned wrong")
}

func assertDate(t *testing.T, ti Date, from string) {
	if !ti.Time.Equal(dateValue) {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, dateValue)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullDate(t *testing.T, ti Date, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
