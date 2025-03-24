package std

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	intJSON = []byte(`12345`)
)

func TestIntFrom(t *testing.T) {
	i := IntFrom(12345)
	assertInt(t, i, "IntFrom()")

	zero := IntFrom(0)
	if !zero.Valid {
		t.Error("IntFrom(0)", "is invalid, but should be valid")
	}
}

func TestIntFromPtr(t *testing.T) {
	n := int64(12345)
	iptr := &n
	i := IntFromPtr(iptr)
	assertInt(t, i, "IntFromPtr()")

	null := IntFromPtr(nil)
	assertNullInt(t, null, "IntFromPtr(nil)")
}

func TestUnmarshalInt(t *testing.T) {
	var i Int
	err := json.Unmarshal(intJSON, &i)
	assert.NoError(t, err)
	assertInt(t, i, "int json")

	var null Int
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assertNullInt(t, null, "null json")

	var badType Int
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt(t, badType, "wrong type json")

	var invalid Int
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}

	assertNullInt(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Int
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	assert.NoError(t, err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt(t *testing.T) {
	var i Int
	err := i.UnmarshalText([]byte("12345"))
	assert.NoError(t, err)
	assertInt(t, i, "UnmarshalText() int")

	var blank Int
	err = blank.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assertNullInt(t, blank, "UnmarshalText() empty int")

	var null Int
	err = null.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assertNullInt(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt(t *testing.T) {
	i := IntFrom(12345)
	data, err := json.Marshal(i)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt(0, false)
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := IntFrom(12345)
	data, err := i.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt(0, false)
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestIntPointer(t *testing.T) {
	i := IntFrom(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewInt(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestIntIsZero(t *testing.T) {
	i := IntFrom(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestIntSetValid(t *testing.T) {
	change := NewInt(0, false)
	assertNullInt(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt(t, change, "SetValid()")
}

func TestIntScan(t *testing.T) {
	var i Int
	err := i.Scan(12345)
	assert.NoError(t, err)
	assertInt(t, i, "scanned int")

	var null Int
	err = null.Scan(nil)
	assert.NoError(t, err)
	assertNullInt(t, null, "scanned null")
}

func TestIntString(t *testing.T) {
	i := IntFrom(12345)
	assert.Equal(t, "12345", i.String())

	null := Int{}
	assert.Equal(t, "", null.String())
}

func TestIntValue(t *testing.T) {
	i := IntFrom(12345)

	v, err := i.Value()
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), v.(int64))

	null := Int{}

	v, err = null.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func assertInt(t *testing.T, i Int, from string) {
	if i.Data != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Data, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt(t *testing.T, i Int, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
