package std

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	uintJSON = []byte(`12345`)
)

func TestUintFrom(t *testing.T) {
	i := UintFrom(12345)
	assertUint(t, i, "UintFrom()")

	zero := UintFrom(0)
	if !zero.Valid {
		t.Error("UintFrom(0)", "is invalid, but should be valid")
	}
}

func TestUintFromPtr(t *testing.T) {
	n := uint64(12345)
	iptr := &n
	i := UintFromPtr(iptr)
	assertUint(t, i, "UintFromPtr()")

	null := UintFromPtr(nil)
	assertNullUint(t, null, "UintFromPtr(nil)")
}

func TestUnmarshalUint(t *testing.T) {
	var i Uint
	err := json.Unmarshal(intJSON, &i)
	assert.NoError(t, err)
	assertUint(t, i, "int json")

	var null Uint
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assertNullUint(t, null, "null json")

	var badType Uint
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint(t, badType, "wrong type json")

	var invalid Uint
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber(t *testing.T) {
	var i Uint
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalUint64Overflow(t *testing.T) {
	uint64Overflow := uint64(math.MaxUint64)

	// Max int64 should decode successfully
	var i Uint
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64Overflow, 10)), &i)
	assert.NoError(t, err)

	// Attempt to overflow
	uint64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64Overflow, 10)), &i)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), i.Data)
}

func TestTextUnmarshalUint(t *testing.T) {
	var i Uint
	err := i.UnmarshalText([]byte("12345"))
	assert.NoError(t, err)
	assertUint(t, i, "UnmarshalText() int")

	var blank Uint
	err = blank.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assertNullUint(t, blank, "UnmarshalText() empty int")

	var null Uint
	err = null.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assertNullUint(t, null, `UnmarshalText() "null"`)
}

func TestMarshalUint(t *testing.T) {
	i := UintFrom(12345)
	data, err := json.Marshal(i)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint(0, false)
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalUintText(t *testing.T) {
	i := UintFrom(12345)
	data, err := i.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint(0, false)
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUintPointer(t *testing.T) {
	i := UintFrom(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewUint(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUintIsZero(t *testing.T) {
	i := UintFrom(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUintSetValid(t *testing.T) {
	change := NewUint(0, false)
	assertNullUint(t, change, "SetValid()")
	change.SetValid(12345)
	assertUint(t, change, "SetValid()")
}

func TestUintScan(t *testing.T) {
	var i Uint
	err := i.Scan(12345)
	assert.NoError(t, err)
	assertUint(t, i, "scanned int")

	var null Uint
	err = null.Scan(nil)
	assert.NoError(t, err)
	assertNullUint(t, null, "scanned null")
}

func TestUintString(t *testing.T) {
	i := UintFrom(12345)
	assert.Equal(t, "12345", i.String())

	null := Uint{}
	assert.Equal(t, "", null.String())
}

func TestUintValue(t *testing.T) {
	i := UintFrom(12345)

	v, err := i.Value()
	assert.NoError(t, err)
	assert.Equal(t, uint64(12345), v.(uint64))

	null := Uint{}

	v, err = null.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func assertUint(t *testing.T, i Uint, from string) {
	if i.Data != uint64(12345) {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Data, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint(t *testing.T, i Uint, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
