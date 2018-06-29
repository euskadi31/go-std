package std

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, dateValue, ti.Data)

	var null Date
	err = json.Unmarshal(nullDateJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var invalid Date
	err = invalid.UnmarshalJSON(invalidJSON)
	assert.Error(t, err)
	assert.False(t, invalid.Valid)

	var bad Date
	err = json.Unmarshal(badObject, &bad)
	assert.Error(t, err)
	assert.False(t, bad.Valid)

	var wrongType Date
	err = json.Unmarshal(intJSON, &wrongType)
	assert.Error(t, err)
	assert.False(t, wrongType.Valid)
}

func TestUnmarshalDateText(t *testing.T) {
	ti := DateFrom(dateValue)
	txt, err := ti.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, dateString, string(txt))

	var unmarshal Date
	err = unmarshal.UnmarshalText(txt)
	assert.NoError(t, err)
	assert.Equal(t, dateValue, unmarshal.Data)

	var null Date
	err = null.UnmarshalText(nullDateJSON)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	txt, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, txt)

	var invalid Date
	err = invalid.UnmarshalText([]byte("hello world"))
	assert.Error(t, err)
	assert.False(t, invalid.Valid)
}

func TestMarshalDate(t *testing.T) {
	dt := Date{}
	data, err := json.Marshal(dt)
	assert.NoError(t, err)
	assert.JSONEq(t, string(nullDateJSON), string(data))

	ti := DateFrom(dateValue)
	data, err = json.Marshal(ti)
	assert.NoError(t, err)
	assert.JSONEq(t, string(dateJSON), string(data))

	ti.Valid = false
	data, err = json.Marshal(ti)
	assert.NoError(t, err)
	assert.JSONEq(t, string(nullDateJSON), string(data))
}

func TestDateFrom(t *testing.T) {
	ti := DateFrom(dateValue)
	assert.True(t, ti.Valid)
	assert.Equal(t, dateValue, ti.Data)
}

func TestDateFromPtr(t *testing.T) {
	ti := DateFromPtr(&dateValue)
	assert.True(t, ti.Valid)
	assert.Equal(t, dateValue, ti.Data)

	null := DateFromPtr(nil)
	assert.False(t, null.Valid)
}

func TestDateSetValid(t *testing.T) {
	var ti time.Time
	change := NewDate(ti, false)
	assert.False(t, change.Valid)

	change.SetValid(dateValue)
	assert.True(t, change.Valid)
	assert.Equal(t, dateValue, change.Data)
}

func TestDatePointer(t *testing.T) {
	ti := DateFrom(dateValue)
	ptr := ti.Ptr()
	assert.Equal(t, dateValue, *ptr)

	var nt time.Time
	null := NewDate(nt, false)
	ptr = null.Ptr()
	assert.Nil(t, ptr)
}

func TestDateScanValue(t *testing.T) {
	var ti Date
	err := ti.Scan(dateValue)
	assert.NoError(t, err)
	assert.True(t, ti.Valid)
	assert.Equal(t, dateValue, ti.Data)

	v, err := ti.Value()
	assert.NoError(t, err)
	assert.Equal(t, dateValue, v)

	var null Date
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	v, err = null.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)

	var wrong Date
	err = wrong.Scan(int64(42))
	assert.Error(t, err)
	assert.False(t, wrong.Valid)
}

func TestDateString(t *testing.T) {
	dt := DateFrom(dateValue)
	assert.Equal(t, "2012-12-21", dt.String())

	null := Date{}
	assert.Equal(t, "", null.String())
}

func TestDateIsZero(t *testing.T) {
	dt := DateFrom(dateValue)
	assert.False(t, dt.IsZero())

	blank := Date{}
	assert.True(t, blank.IsZero())

	empty := NewDate(time.Time{}, true)
	assert.False(t, empty.IsZero())

	null := DateFromPtr(nil)
	assert.True(t, null.IsZero())
}
