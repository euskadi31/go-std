package std

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)

	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)
)

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assert.True(t, str.Valid)
	assert.Equal(t, "test", str.Data)

	zero := StringFrom("")
	assert.True(t, zero.Valid)
	assert.Equal(t, "", zero.Data)
}

func TestStringFromPtr(t *testing.T) {
	s := "test"
	sptr := &s
	str := StringFromPtr(sptr)
	assert.True(t, str.Valid)
	assert.Equal(t, "test", str.Data)

	null := StringFromPtr(nil)
	assert.False(t, null.Valid)
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	assert.NoError(t, err)
	assert.True(t, str.Valid)
	assert.Equal(t, "test", str.Data)

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	assert.NoError(t, err)
	assert.True(t, blank.Valid)
	assert.Equal(t, "", blank.Data)

	var null String
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var badType String
	err = json.Unmarshal(boolJSON, &badType)
	assert.Error(t, err)
	assert.False(t, badType.Valid)

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	assert.Error(t, err)
	assert.False(t, invalid.Valid)
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	assert.NoError(t, err)
	assert.True(t, str.Valid)
	assert.Equal(t, "test", str.Data)

	var null String
	err = null.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func TestMarshalString(t *testing.T) {
	str := StringFrom("test")
	data, err := json.Marshal(str)
	assert.NoError(t, err)
	assert.JSONEq(t, `"test"`, string(data))

	data, err = str.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))

	// empty values should be encoded as an empty string
	zero := StringFrom("")
	data, err = json.Marshal(zero)
	assert.NoError(t, err)
	assert.JSONEq(t, `""`, string(data))

	data, err = zero.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))

	null := StringFromPtr(nil)
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.JSONEq(t, `null`, string(data))

	data, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func TestStringPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Ptr()
	assert.Equal(t, "test", *ptr)

	null := NewString("", false)
	ptr = null.Ptr()
	assert.Nil(t, ptr)
}

func TestStringIsZero(t *testing.T) {
	str := StringFrom("test")
	assert.False(t, str.IsZero())

	blank := StringFrom("")
	assert.False(t, blank.IsZero())

	empty := NewString("", true)
	assert.False(t, empty.IsZero())

	null := StringFromPtr(nil)
	assert.True(t, null.IsZero())
}

func TestStringSetValid(t *testing.T) {
	change := NewString("", false)
	assert.False(t, change.Valid)

	change.SetValid("test")
	assert.True(t, change.Valid)
	assert.Equal(t, "test", change.Data)
}

func TestStringScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	assert.NoError(t, err)
	assert.True(t, str.Valid)
	assert.Equal(t, "test", str.Data)

	var null String
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func TestStringString(t *testing.T) {
	str := StringFrom("test")
	assert.Equal(t, "test", str.String())

	null := String{}
	assert.Equal(t, "", null.String())
}

func TestStringValue(t *testing.T) {
	str := StringFrom("test")

	v, err := str.Value()
	assert.NoError(t, err)
	assert.Equal(t, "test", v.(string))

	null := String{}

	v, err = null.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
