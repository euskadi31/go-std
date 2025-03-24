package std

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	boolJSON = []byte(`true`)
)

func TestBoolFrom(t *testing.T) {
	b := BoolFrom(true)

	assert.True(t, b.Valid)
	assert.True(t, b.Data)

	zero := BoolFrom(false)

	assert.True(t, zero.Valid)
	assert.False(t, zero.Data)
}

func TestBoolFromPtr(t *testing.T) {
	n := true
	bptr := &n
	b := BoolFromPtr(bptr)

	assert.True(t, b.Valid)
	assert.True(t, b.Data)

	null := BoolFromPtr(nil)

	assert.False(t, null.Valid)
}

func TestUnmarshalBool(t *testing.T) {
	var b Bool
	err := json.Unmarshal(boolJSON, &b)
	assert.NoError(t, err)
	assert.True(t, b.Valid)
	assert.True(t, b.Data)

	var null Bool
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var badType Bool
	err = json.Unmarshal(intJSON, &badType)
	assert.Error(t, err)
	assert.False(t, badType.Valid)

	var invalid Bool
	err = invalid.UnmarshalJSON(invalidJSON)
	assert.EqualError(t, err, "json: cannot unmarshal :) into Go value of type null.Bool: invalid character ':' looking for beginning of value")
}

func TestTextUnmarshalBool(t *testing.T) {
	var b Bool
	err := b.UnmarshalText(boolJSON)
	assert.NoError(t, err)
	assert.True(t, b.Valid)
	assert.True(t, b.Data)

	var zero Bool
	err = zero.UnmarshalText([]byte("false"))
	assert.NoError(t, err)
	assert.True(t, zero.Valid)
	assert.False(t, zero.Data)

	var blank Bool
	err = blank.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assert.False(t, blank.Valid)

	var null Bool
	err = null.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var invalid Bool
	err = invalid.UnmarshalText([]byte(":D"))
	assert.Error(t, err)
	assert.False(t, null.Valid)
}

func TestMarshalBool(t *testing.T) {
	b := BoolFrom(true)
	data, err := json.Marshal(b)
	assert.NoError(t, err)
	assert.JSONEq(t, `true`, string(data))

	zero := NewBool(false, true)
	data, err = json.Marshal(zero)
	assert.NoError(t, err)
	assert.JSONEq(t, `false`, string(data))

	// invalid values should be encoded as null
	null := NewBool(false, false)
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.JSONEq(t, `null`, string(data))
}

func TestMarshalBoolText(t *testing.T) {
	b := BoolFrom(true)
	data, err := b.MarshalText()
	assert.NoError(t, err)
	assert.JSONEq(t, `true`, string(data))

	zero := NewBool(false, true)
	data, err = zero.MarshalText()
	assert.NoError(t, err)
	assert.JSONEq(t, `false`, string(data))

	// invalid values should be encoded as null
	null := NewBool(false, false)
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func TestBoolPointer(t *testing.T) {
	b := BoolFrom(true)
	ptr := b.Ptr()
	assert.True(t, *ptr)

	null := NewBool(false, false)
	ptr = null.Ptr()
	assert.Nil(t, ptr)
}

func TestBoolIsZero(t *testing.T) {
	b := BoolFrom(true)
	assert.False(t, b.IsZero())

	null := NewBool(false, false)
	assert.True(t, null.IsZero())

	zero := NewBool(false, true)
	assert.False(t, zero.IsZero())
}

func TestBoolSetValid(t *testing.T) {
	change := NewBool(false, false)
	assert.False(t, change.Valid)

	change.SetValid(true)
	assert.True(t, change.Valid)
}

func TestBoolScan(t *testing.T) {
	var b Bool
	err := b.Scan(true)
	assert.NoError(t, err)
	assert.True(t, b.Valid)
	assert.True(t, b.Data)

	var null Bool
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func TestBoolString(t *testing.T) {
	b := BoolFrom(true)
	assert.Equal(t, "true", b.String())

	b = BoolFrom(false)
	assert.Equal(t, "false", b.String())

	b = Bool{}
	assert.Equal(t, "", b.String())
}

func TestBoolValue(t *testing.T) {
	b := BoolFrom(true)

	v, err := b.Value()
	assert.NoError(t, err)
	assert.True(t, v.(bool))

	b = Bool{}

	v, err = b.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}
