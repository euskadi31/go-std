package std

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	dateTimeString    = "2012-12-21T21:21:21+0000"
	dateTimeJSON      = []byte(`"` + dateTimeString + `"`)
	nullDateTimeJSON  = []byte(`null`)
	dateTimeValue, _  = time.Parse(dateTimeFormat, dateTimeString)
	badDateTimeObject = []byte(`{"hello": "world"}`)
)

func TestUnmarshalDateTimeJSON(t *testing.T) {
	var ti DateTime
	err := json.Unmarshal(dateTimeJSON, &ti)
	assert.NoError(t, err)
	assertDateTime(t, ti, "UnmarshalJSON() json")

	var null DateTime
	err = json.Unmarshal(nullDateTimeJSON, &null)
	assert.NoError(t, err)
	assertNullDateTime(t, null, "null time json")

	var invalid DateTime
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullDateTime(t, invalid, "invalid from object json")

	var bad DateTime
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullDateTime(t, bad, "bad from object json")

	var wrongType DateTime
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullDateTime(t, wrongType, "wrong type object json")
}

func TestUnmarshalDateTimeText(t *testing.T) {
	ti := DateTimeFrom(dateTimeValue)
	txt, err := ti.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, txt, dateTimeString, "marshal text")

	var unmarshal DateTime
	err = unmarshal.UnmarshalText(txt)
	assert.NoError(t, err)
	assertDateTime(t, unmarshal, "unmarshal text")

	var null DateTime
	err = null.UnmarshalText(nullDateTimeJSON)
	assert.NoError(t, err)
	assertNullDateTime(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, txt)

	var invalid DateTime
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDateTime(t, invalid, "bad string")
}

func TestMarshalDateTime(t *testing.T) {
	dt := DateTime{}
	data, err := json.Marshal(dt)
	assert.NoError(t, err)
	assert.JSONEq(t, `null`, string(data))

	ti := DateTimeFrom(dateTimeValue)
	data, err = json.Marshal(ti)
	assert.NoError(t, err)
	assert.JSONEq(t, string(dateTimeJSON), string(data))

	ti.Valid = false
	data, err = json.Marshal(ti)
	assert.NoError(t, err)
	assert.JSONEq(t, `null`, string(data))
}

func TestDateTimeFrom(t *testing.T) {
	ti := DateTimeFrom(dateTimeValue)
	assertDateTime(t, ti, "DateTimeFrom() time.Time")
}

func TestDateTimeFromPtr(t *testing.T) {
	ti := DateTimeFromPtr(&dateTimeValue)
	assertDateTime(t, ti, "DateTimeFromPtr() time")

	null := DateTimeFromPtr(nil)
	assertNullDateTime(t, null, "DateTimeFromPtr(nil)")
}

func TestDateTimeSetValid(t *testing.T) {
	var ti time.Time
	change := NewDateTime(ti, false)
	assertNullDateTime(t, change, "SetValid()")
	change.SetValid(dateTimeValue)
	assertDateTime(t, change, "SetValid()")
}

func TestDateTimePointer(t *testing.T) {
	ti := DateTimeFrom(dateTimeValue)
	ptr := ti.Ptr()
	if *ptr != dateTimeValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, dateTimeValue)
	}

	var nt time.Time
	null := NewDateTime(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestDateTimeScanValue(t *testing.T) {
	var ti DateTime
	err := ti.Scan(dateTimeValue)
	assert.NoError(t, err)
	assertDateTime(t, ti, "scanned time")
	if v, err := ti.Value(); v != dateTimeValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null DateTime
	err = null.Scan(nil)
	assert.NoError(t, err)
	assertNullDateTime(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong DateTime
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDateTime(t, wrong, "scanned wrong")
}

func TestDateTimeString(t *testing.T) {
	dt := DateTimeFrom(dateTimeValue)
	assert.Equal(t, "2012-12-21T21:21:21+0000", dt.String())

	null := DateTime{}
	assert.Equal(t, "", null.String())
}

func TestDateTimeIsZero(t *testing.T) {
	dt := DateTimeFrom(dateTimeValue)
	assert.False(t, dt.IsZero())

	blank := DateTime{}
	assert.True(t, blank.IsZero())

	empty := NewDateTime(time.Time{}, true)
	assert.False(t, empty.IsZero())

	null := DateTimeFromPtr(nil)
	assert.True(t, null.IsZero())
}

func assertDateTime(t *testing.T, ti DateTime, from string) {
	if !ti.Data.Equal(dateTimeValue) {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Data, dateTimeValue)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullDateTime(t *testing.T, ti DateTime, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
